export function Loading() {
  const spinnerStyle = {
    width: "30px",
    height: "30px",
    border: "4px solid rgba(0, 0, 0, 0.1)",
    borderLeftColor: "#3498db",
    borderRadius: "50%",
    animation: "spin 1s linear infinite",
  };

  return (
    <div
      className="spinner-container"
      style={{ display: "flex", justifyContent: "center", padding: "20px" }}
    >
      <div className="loading-spinner" style={spinnerStyle}>
        <style>
          {`
            @keyframes spin {
              0% { transform: rotate(0deg); }
              100% { transform: rotate(360deg); }
            }
          `}
        </style>
      </div>
    </div>
  );
}
