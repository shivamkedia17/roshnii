import dummyProfilePhotoPath from "@/dummyData/profile.png";

export function Profile() {
  return (
    <div className="profile">
      <div className="profile-photo">
        <img src={dummyProfilePhotoPath} alt="User profile" />
      </div>
      <h1>Full Name</h1>
      <p>Other Info As Needed</p>
    </div>
  );
}
