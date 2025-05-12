import { useGetCurrentUser } from "@/hooks/useUser";
import { useState } from "react";
import "@/css/Profile.css";

export function Profile() {
  const { currentUser, isLoading, error } = useGetCurrentUser();
  const [imageError, setImageError] = useState(false);

  if (isLoading) {
    return (
      <div className="profile-container loading">
        <div className="loading-spinner"></div>
        <p>Loading profile...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="profile-container error">
        <h2>Error Loading Profile</h2>
        <p>{error.message || "An unexpected error occurred"}</p>
      </div>
    );
  }

  if (!currentUser) {
    return (
      <div className="profile-container error">
        <h2>User Not Found</h2>
        <p>Unable to retrieve user profile information.</p>
      </div>
    );
  }

  return (
    <div className="profile-container">
      <div className="profile-header">
        <div className="profile-photo">
          {currentUser.picture_url && !imageError ? (
            <img
              src={currentUser.picture_url}
              alt={`${currentUser.name}'s profile`}
              onError={() => setImageError(true)}
              className="user-avatar"
            />
          ) : (
            <div className="user-avatar placeholder">
              {currentUser.name
                ? currentUser.name.charAt(0).toUpperCase()
                : "U"}
            </div>
          )}
        </div>
        <div className="profile-info">
          <h1>{currentUser.name || "User"}</h1>
          <p className="user-email">{currentUser.email}</p>
        </div>
      </div>

      <div className="profile-details">
        <div className="detail-card">
          <h3>Account Information</h3>
          <div className="detail-row">
            <span className="detail-label">Authentication:</span>
            <span className="detail-value">Google OAuth</span>
          </div>
        </div>
      </div>
    </div>
  );
}
