import React, { useState, useEffect } from "react";
import { Box, Table, TableBody, TableCell as StyledTableCell, TableContainer, TableHead, TableRow, Paper, ThemeProvider, createTheme } from "@mui/material";
import LockIcon from '@mui/icons-material/Lock';
import LockOpenIcon from '@mui/icons-material/LockOpen';

const App = () => {
  const [racks, setRacks] = useState([]);

  useEffect(() => {
    fetch("http://localhost:8080/api/racks")
      .then(async (response) => {
        const text = await response.text();
        try {
          return JSON.parse(text);
        } catch (error) {
          throw new Error(`Invalid JSON response: ${text}`);
        }
      })
      .then((data) => {
        console.log("Racks:", data.rack_details);
        setRacks(data.rack_details);
      })
      .catch((error) => console.error("Fetch error:", error));

    const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = () => {
      console.log("WebSocket connected");
    };

    ws.onmessage = (event) => {
      console.log("Received WebSocket message:", event.data);
      const updatedRack = JSON.parse(event.data);

      setRacks((prevRacks) =>
        prevRacks.map((rack) =>
          rack.rack_id === updatedRack.rack_id
            ? {
                ...rack,
                user_id: updatedRack.user_id,
                user_name: updatedRack.user_name,
                user_email: updatedRack.user_email,
                user_phone: updatedRack.user_phone,
                last_updated: updatedRack.last_updated,
              }
            : rack
        )
      );
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    ws.onclose = () => {
      console.log("WebSocket disconnected");
    };

    return () => ws.close();
  }, []);

  const getCellStyles = (userId) => ({
    color: userId !== 0 ? "white" : "#D3D3D3",
    fontFamily: 'Courier New',
    fontWeight: userId !== 0 ? "bold" : "normal",
    backgroundColor: userId !== 0 ? "#555" : "transparent"
  });

  return (
    <Box sx={{ width: "80%", marginTop: 2, marginLeft: "auto", marginRight: "auto", borderRadius: 10, display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center", minHeight: "100vh" }}>
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="rack table">
          <TableHead sx={{ backgroundColor: "#353935" }}>
            <TableRow>
              {["Rack ID", "User ID", "User Name", "User Email", "User Phone", "Time Locked", "Status", "Action"].map((label) => (
                <StyledTableCell align="center" sx={{ color: "#808080", fontWeight: "bolder", fontFamily: 'Courier New' }} key={label}>
                  {label}
                </StyledTableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody sx={{ backgroundColor: "#353935" }}>
            {racks.map((rack) => (
              <TableRow key={rack.rack_id} sx={{ backgroundColor: rack.user_id !== 0 ? "#333" : "transparent" }}>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>
                  {rack.rack_id}
                </StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>
                  {rack.user_id !== 0 ? rack.user_id : null}
                </StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>
                  {rack.user_name}
                </StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>
                  {rack.user_email}
                </StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>
                  {rack.user_phone}
                </StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>
                  {rack.last_updated ? new Date(rack.last_updated).toLocaleString('en-US', { hour12: false, timeZone: 'UTC' }).slice(0, 16) : ""}
                </StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>
                  {rack.user_id !== 0 ? (
                    <LockIcon sx={{ color: "#FFCCCB", fontSize: 30 }} />
                  ) : (
                    <LockOpenIcon sx={{ color: "#66FF99", fontSize: 30 }} />
                  )}
                </StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>
                  <button
                    style={{
                      backgroundColor: rack.user_id === 0 ? "#A9A9A9" : "#CF9FFF",
                      color: "white",
                      border: "none",
                      padding: "8px 16px",
                      borderRadius: "4px",
                      cursor: rack.user_id === 0 ? "not-allowed" : "pointer",
                      fontFamily: 'Courier New',
                      fontWeight: "bold",
                    }}
                    onClick={() => {
                      if (rack.user_id !== 0) {
                        console.log("Unlock Rack:", rack.rack_id);
                      }
                    }}
                    disabled={rack.user_id === 0}
                  >
                    UNLOCK
                  </button>
                </StyledTableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  );
};

export default App;
