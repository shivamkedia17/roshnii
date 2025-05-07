import "@/dummyData/profile.png";

export function Profile() {
  return (
    <div className="profile">
      <div className="profile-photo">
        <img src="@/dummyData/profile.png" alt="User profile" />
      </div>
      <h1>Full Name</h1>
      <p>Other Info As Needed</p>
    </div>
  );
}
