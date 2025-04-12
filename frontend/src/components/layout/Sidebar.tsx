import "../../css/Sidebar.css";

interface SidebarProps {
  isOpen: boolean;
  activeView: string;
  onNavigate: (view: string) => void;
}

export function Sidebar({ isOpen, activeView, onNavigate }: SidebarProps) {
  return (
    <aside className={`sidebar ${isOpen ? "open" : "closed"}`}>
      <nav className="sidebar-nav">
        <ul>
          <li className={activeView === "photos" ? "active" : ""}>
            <button onClick={() => onNavigate("photos")}>
              <span className="icon">📷</span>
              <span className="label">Photos</span>
            </button>
          </li>
          <li className={activeView === "albums" ? "active" : ""}>
            <button onClick={() => onNavigate("albums")}>
              <span className="icon">📁</span>
              <span className="label">Albums</span>
            </button>
          </li>
          <li className={activeView === "upload" ? "active" : ""}>
            <button onClick={() => onNavigate("upload")}>
              <span className="icon">⬆️</span>
              <span className="label">Upload</span>
            </button>
          </li>
        </ul>
      </nav>
    </aside>
  );
}
