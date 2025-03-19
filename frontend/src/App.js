import React, { useState, useEffect } from "react";
import { Box, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, IconButton, TableCell as StyledTableCell, CircularProgress, TextField, Tooltip } from "@mui/material";
import LockIcon from '@mui/icons-material/Lock';
import LockOpenIcon from '@mui/icons-material/LockOpen';
import FilterListIcon from '@mui/icons-material/FilterList';
import logo from './assets/logo.png';

const App = () => {
  const [racks, setRacks] = useState([]);
  const [loadingRack, setLoadingRack] = useState(-1);
  const [filters, setFilters] = useState({});
  const [filteringColumn, setFilteringColumn] = useState(null);

  useEffect(() => {
    fetch("http://localhost:8080/api/racks")
      .then((response) => response.json())
      .then((data) => setRacks(data.rack_details))
      .catch((error) => console.error("Fetch error:", error));

    const ws = new WebSocket("ws://localhost:8080/ws");
    ws.onmessage = (event) => {
      const updatedRack = JSON.parse(event.data);
      setRacks((prevRacks) => prevRacks.map((rack) => (rack.rack_id === updatedRack.rack_id ? updatedRack : rack)));
    };
    return () => ws.close();
  }, []);

  const toggleFilter = (key) => {
    // Toggle the filter visibility for the clicked column
    setFilteringColumn(filteringColumn === key ? null : key);
  };

  const handleFilterChange = (value, key) => {
    const filterValue = value.toLowerCase();
    setFilters((prevFilters) => {
      if (filterValue === "") {
        // Remove the filter if it's cleared
        const { [key]: _, ...rest } = prevFilters;
        return rest;
      }
      return { ...prevFilters, [key]: filterValue };
    });
  };

  const filteredRacks = racks.filter((rack) =>
    Object.keys(filters).every((key) => {
      let rack_key = key.toLowerCase().replace(" ", "_");
      return rack[rack_key]?.toString().toLowerCase().includes(filters[key] || "");
    })
  );

  const unlockRack = (rackId) => {
    setLoadingRack(rackId);
    fetch("http://localhost:8080/api/unlock_frontend", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ rack_id: rackId })
    })
    .finally(() => {
      setTimeout(() => {setLoadingRack(-1);}, 1000);
    });
  };

  const getCellStyles = (userId) => ({
    color: userId !== 0 ? "white" : "#D3D3D3",
    fontFamily: 'Courier New',
    fontWeight: userId !== 0 ? "bold" : "normal",
    backgroundColor: userId !== 0 ? "#555" : "transparent"
  });

  return (
    <Box sx={{ width: "100%", display: "flex", flexDirection: "column", alignItems: "center", minHeight: "100vh", backgroundColor: "#353935" }}>
      <Box sx={{ display: "flex", alignItems: "center", gap: 2, marginBottom: 2, marginTop: 2 }}>
        <img src={logo} alt="Logo" style={{ height: 50 }} />
        <h1 style={{ color: "#CF9FFF", fontFamily: 'Courier New', fontWeight: "bolder" }}>Anchor Dashboard</h1>
      </Box>
      <TableContainer component={Paper} sx={{ width: "80%", backgroundColor: "#353935", borderColor: "#CF9FFF", boxShadow: "0 0 10px #CF9FFF" }}>
        <Table sx={{ minWidth: 650 }}>
          <TableHead>
            <TableRow>
              {["Rack ID", "Location", "User ID", "User Name", "User Email", "User Phone", "Time Locked", "Status", "Action"].map((label, index) => (
                <TableCell key={index} align="center" sx={{ color: "#808080", fontWeight: "bolder" }}>
                  {label}
                  {index < 6 && (
                    <Tooltip
                      title={
                        filteringColumn === label && (
                          <TextField
                            size="small"
                            variant="outlined"
                            placeholder={`Filter ${label}`}
                            value={filters[label] || ""}
                            onChange={(e) => handleFilterChange(e.target.value, label)}
                            autoFocus
                            onBlur={() => setFilteringColumn(null)}  // Close the tooltip when it loses focus
                          />
                        )
                      }
                      open={filteringColumn === label}
                      onClose={() => setFilteringColumn(null)}
                      disableFocusListener
                      disableHoverListener
                      disableTouchListener
                    >
                      <IconButton onClick={() => toggleFilter(label)}>
                        <FilterListIcon sx={{ color: filters[label] ? "#CF9FFF" : "#808080" }} />
                      </IconButton>
                    </Tooltip>
                  )}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {filteredRacks.map((rack) => (
              <TableRow key={rack.rack_id} sx={{ backgroundColor: rack.user_id !== 0 ? "#333" : "transparent" }}>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>{rack.rack_id}</StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>Engineering 7 Building</StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>{rack.user_id || ""}</StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>{rack.user_name}</StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>{rack.user_email}</StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>{rack.user_phone}</StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>{rack.last_updated ? new Date(rack.last_updated).toLocaleString() : ""}</StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>
                  {rack.user_id !== 0 ? <LockIcon sx={{ color: "#FFCCCB" }} /> : <LockOpenIcon sx={{ color: "#66FF99" }} />}
                </StyledTableCell>
                <StyledTableCell align="center" sx={getCellStyles(rack.user_id)}>
                  {loadingRack !== rack.rack_id ? (
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
                      onClick={() => unlockRack(rack.rack_id)}
                      disabled={rack.user_id === 0}
                    >
                      UNLOCK
                    </button>
                  ) : (
                    <CircularProgress size={24} sx={{ color: "#CF9FFF" }} />
                  )}
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
