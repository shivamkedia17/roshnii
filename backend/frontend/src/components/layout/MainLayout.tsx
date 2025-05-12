import { useState, useEffect } from "react";
import { Gallery } from "../photos/PhotoGallery";
import { Header } from "./Header";
import { Sidebar } from "./Sidebar";
import { UploadForm } from "../upload/UploadForm";
import { Profile } from "./Profile";
import { AlbumGallery } from "../albums/AlbumGallery";

import "@/css/MainLayout.css";
import "@/css/AlbumView.css";

export function MainLayout() {
  const [activeView, setActiveView] = useState("photos");
  const [selectedAlbumId, setSelectedAlbumId] = useState<string | null>(null);
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [_, setSearchQuery] = useState("");

  // Handle album selection - updates albumId and switches view if needed
  const handleSelectAlbum = (albumId: string) => {
    setSelectedAlbumId(albumId);
    if (activeView !== "albums") {
      setActiveView("albums");
    }
  };

  // Clear album selection when navigating away from albums
  useEffect(() => {
    if (activeView !== "albums") {
      setSelectedAlbumId(null);
    }
  }, [activeView]);

  function renderMainContent() {
    switch (activeView) {
      case "photos":
        return <Gallery albumId={undefined} />;
      case "albums":
        // If an album is selected, show Gallery with that albumId, otherwise show the AlbumGallery
        return selectedAlbumId ? (
          <div className="album-view">
            <div className="album-view-header">
              <button
                className="back-button"
                onClick={() => setSelectedAlbumId(null)}
              >
                ← Back to Albums
              </button>
            </div>
            <div className="album-content">
              <Gallery albumId={selectedAlbumId} />
            </div>
          </div>
        ) : (
          <AlbumGallery onSelectAlbum={handleSelectAlbum} />
        );
      case "upload":
        return <UploadForm onComplete={() => setActiveView("photos")} />;

      case "profile":
        return <Profile />;
      case "faces":
        return <div>Faces feature coming soon</div>;
      default:
        return <Gallery />;
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
          onNavigate={(view) => {
            setActiveView(view);
          }}
        />
        {sidebarOpen && (
          <div
            className={`sidebar-overlay ${sidebarOpen ? "visible" : ""}`}
            onClick={() => setSidebarOpen(false)}
          />
        )}
        <main className="main-content">{renderMainContent()}</main>
      </div>
    </div>
  );
}
