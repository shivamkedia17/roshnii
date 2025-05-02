import { useState } from "react";
import { Header } from "./Header";
import { Sidebar } from "./Sidebar";
import { Gallery } from "../photos/Gallery";
import { AlbumList } from "../albums/AlbumList";
import { UploadForm } from "../upload/UploadForm";
import { AlbumView } from "../albums/AlbumView";

import "@/css/MainLayout.css";

export function MainLayout() {
  const [activeView, setActiveView] = useState("photos");
  const [selectedAlbumId, setSelectedAlbumId] = useState<string | null>(null);
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");

  function renderMainContent() {
    switch (activeView) {
      case "photos": // TODO
        return <Gallery searchQuery={searchQuery} />;
      case "albums": // TODO
        return selectedAlbumId ? (
          <AlbumView
            albumId={selectedAlbumId}
            onBack={() => setSelectedAlbumId(null)}
          />
        ) : (
          <AlbumList onSelectAlbum={setSelectedAlbumId} />
        );
      case "upload": // TODO
        return <UploadForm onComplete={() => setActiveView("photos")} />;
      case "search": // TODO
        return <Gallery searchQuery={searchQuery} />;
      case "profile": // TODO
        return <></>;
      case "logout": // TODO
        return <></>;

      default:
        return <Gallery searchQuery={searchQuery} />;
    }
  }

  return (
    <div className="main-layout">
      <Header
        onToggleSidebar={() => setSidebarOpen(!sidebarOpen)}
        onSearch={setSearchQuery}
      />
      <div className="content-container">
        <Sidebar
          isOpen={sidebarOpen}
          activeView={activeView}
          onNavigate={setActiveView}
        />
        {/* Only show overlay on mobile */}
        {sidebarOpen && (
          <div
            className={`sidebar-overlay ${sidebarOpen ? "visible" : ""}`}
            onClick={() => setSidebarOpen(false)}
          />
        )}
        <main className={`main-content `}>{renderMainContent()}</main>
      </div>
    </div>
  );
}
