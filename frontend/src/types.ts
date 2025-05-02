export type HeaderProps = {
  onToggleSidebar: () => void; // function that mutates state
  onSearch: (query: string) => void; // similar
};

export type SidebarProps = {
  isOpen: boolean;
  activeView: string;
  onNavigate: (view: string) => void;
};

export type GalleryProps = {
  searchQuery?: string;
  albumId?: string;
};
