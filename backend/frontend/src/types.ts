export type HeaderProps = {
  onToggleSidebar: () => void;
  onSearch: (query: string) => void;
};

export type SidebarProps = {
  isOpen: boolean;
  activeView: string;
  onNavigate: (view: string) => void;
};

export type PhotoModalProps = {
  photoId: string;
  onClose: () => void;
  albumContext?: {
    albumId: number;
    onRemoveFromAlbum: (photoId: string) => void;
  };
};
