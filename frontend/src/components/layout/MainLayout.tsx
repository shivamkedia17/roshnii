import { useState } from "react";
import { Header } from "./Header";
import { Sidebar } from "./Sidebar";
import { Gallery } from "../photos/Gallery";
import { AlbumList } from "../albums/AlbumList";
import { AlbumView } from "../albums/AlbumView"; // Import the AlbumView component
import { UploadForm } from "../upload/UploadForm";
import { Profile } from "./Profile";
import { LogoutPage } from "../auth/LogoutPage";

import "@/css/MainLayout.css";

export function MainLayout() {
  const [activeView, setActiveView] = useState("photos");
  const [selectedAlbumId, setSelectedAlbumId] = useState<string | null>(null);
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");

  function renderMainContent() {
    switch (activeView) {
      case "photos":
        return (
          <Gallery
            searchQuery={searchQuery}
            albumId={selectedAlbumId ?? undefined}
          />
        );
      case "albums":
        // If an album is selected, show the AlbumView, otherwise show the AlbumList
        return selectedAlbumId ? (
          <AlbumView
            albumId={selectedAlbumId}
            onBack={() => setSelectedAlbumId(null)}
          />
        ) : (
          <AlbumList onSelectAlbum={setSelectedAlbumId} />
        );
      case "upload":
        return <UploadForm onComplete={() => setActiveView("photos")} />;
      case "profile":
        return <Profile />;
      case "faces":
        return <div>Faces feature coming soon</div>;
      case "logout":
        return <LogoutPage />;
      default:
        return <Gallery searchQuery={searchQuery} />;
    }
  }

  // Handle navigating to albums view and selecting an album
  const handleNavigateToAlbum = (albumId: string | null) => {
    setSelectedAlbumId(albumId);
    setActiveView("albums");
  };

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
            // Reset selected album when navigating away from albums
            if (view !== "albums") {
              setSelectedAlbumId(null);
            }
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
