import { SidebarProps } from "@/types";
import { useLogout } from "@/hooks/useAuth";
import "../../css/Sidebar.css";

export function Sidebar({ isOpen, activeView, onNavigate }: SidebarProps) {
  const logout = useLogout();

  return (
    <aside className={`sidebar ${isOpen ? "open" : "closed"}`}>
      <nav className="sidebar-nav">
        <button
          className={activeView === "upload" ? "active" : ""}
          onClick={() => onNavigate("upload")}
        >
          <span className="icon">🧧</span>
          <span className="label">Upload</span>
        </button>
        <button
          className={activeView === "photos" ? "active" : ""}
          onClick={() => onNavigate("photos")}
        >
          <span className="icon">📷</span>
          <span className="label">Photos</span>
        </button>
        <button
          className={activeView === "albums" ? "active" : ""}
          onClick={() => onNavigate("albums")}
        >
          <span className="icon">🏞️</span>
          <span className="label">Albums</span>
        </button>
        <button
          className={activeView === "faces" ? "active" : ""}
          onClick={() => onNavigate("faces")}
        >
          <span className="icon">🤦‍♀️</span>
          <span className="label">Faces</span>
        </button>
      </nav>

      <nav className="sidebar-userinfo">
        <button
          className={activeView === "profile" ? "active" : ""}
          onClick={() => onNavigate("profile")}
        >
          <span className="icon">👤</span>
          <span className="label">Profile</span>
        </button>
        <button
          className={activeView === "logout" ? "active" : ""}
          onClick={() => logout()}
        >
          <span className="label">Logout</span>
        </button>
      </nav>
    </aside>
  );
}
