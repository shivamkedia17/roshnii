import { useState } from "react";
import { Header } from "./Header";
import { Sidebar } from "./Sidebar";
import { Gallery } from "../photos/Gallery";
import { AlbumList } from "../albums/AlbumList";
// import { AlbumView } from "../albums/AlbumView";
import { UploadForm } from "../upload/UploadForm";
import "../../css/Layout.css";

export function MainLayout() {
  const [activeView, setActiveView] = useState("photos");
  const [selectedAlbumId, setSelectedAlbumId] = useState<string | null>(null);
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");

  const renderMainContent = () => {
    switch (activeView) {
      case "photos":
        return <Gallery searchQuery={searchQuery} />;
      case "albums":
        return selectedAlbumId ? (
          <AlbumDetail
            albumId={selectedAlbumId}
            onBack={() => setSelectedAlbumId(null)}
          />
        ) : (
          <AlbumList onSelectAlbum={setSelectedAlbumId} />
        );
      case "upload":
        return <UploadForm onComplete={() => setActiveView("photos")} />;
      case "search":
        return <Gallery searchQuery={searchQuery} />;
      default:
        return <Gallery searchQuery={searchQuery} />;
    }
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
          onNavigate={setActiveView}
        />
        <main className={`main-content ${sidebarOpen ? "sidebar-open" : ""}`}>
          {renderMainContent()}
        </main>
      </div>
    </div>
  );
}
