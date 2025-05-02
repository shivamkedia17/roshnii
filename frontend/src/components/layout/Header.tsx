import { useState } from "react";
import { HeaderProps } from "@/types";
import "../../css/Header.css";

export function Header({ onToggleSidebar, onSearch }: HeaderProps) {
  const [searchQuery, setSearchQuery] = useState("");

  function handleSearch(e: React.FormEvent) {
    e.preventDefault();
    onSearch(searchQuery);
  }

  return (
    <header className="app-header">
      <div className="header-left">
        <h3 className="app-title">Roshnii</h3>
        <button className="menu-toggle" onClick={onToggleSidebar}>
          <span className="menu-icon">
            <svg
              viewBox="0 0 24 24"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
              stroke="currentColor"
              strokeWidth="1.2"
              strokeLinecap="round"
              className="hamburger"
            >
              <g id="SVGRepo_iconCarrier">
                <path d="M20 7L4 7"></path>
                <path d="M20 12L4 12"></path>
                <path d="M20 17L4 17"></path>
              </g>
            </svg>
          </span>
          <span>Menu</span>
        </button>
      </div>
      <div className="header-center">
        <form className="search-form" onSubmit={handleSearch}>
          <input
            type="text"
            placeholder="Search photos..."
            value={searchQuery}
            className="search-bar"
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          <button type="submit">🔍</button>
        </form>
      </div>
    </header>
  );
}
