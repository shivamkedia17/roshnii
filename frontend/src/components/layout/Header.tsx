import { useState } from "react";
import "../../css/Header.css";

interface HeaderProps {
  onToggleSidebar: () => void;
  onSearch: (query: string) => void;
}

export function Header({ onToggleSidebar, onSearch }: HeaderProps) {
  const [searchQuery, setSearchQuery] = useState("");

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    onSearch(searchQuery);
  };

  return (
    <header className="app-header">
      <div className="header-left">
        <button className="menu-toggle" onClick={onToggleSidebar}>
          <span className="menu-icon">☰</span>
        </button>
        <h1 className="app-title">Roshnii Photos</h1>
      </div>

      <form className="search-form" onSubmit={handleSearch}>
        <input
          type="text"
          placeholder="Search photos..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
        <button type="submit">Search</button>
      </form>

      <div className="user-menu">
        <img
          className="user-avatar"
          src="/profile-placeholder.jpg"
          alt="Profile"
        />
      </div>
    </header>
  );
}
