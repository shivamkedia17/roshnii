export type AuthTokens = {
  // Simplify to focus on cookie-based auth response
  message: string;
  // Only used in development mode
  token?: string;
  expiresIn?: number;
};

export type AuthContextType = {
  isAuthenticated: boolean;
  isLoading: boolean;
  user: UserInfo | null;
  login: () => void;
  logout: () => Promise<boolean | void>;
  refreshToken: () => Promise<boolean>;
};

export type UserInfo = {
  user_id: number;
  email: string;
  name?: string;
  picture_url?: string;
};

export type PhotoInfo = {
  id: string;
  user_id: number;
  filename: string;
  storage_path?: string;
  content_type: string;
  size: number;
  width: number | null;
  height: number | null;
  created_at: string;
  updated_at: string;
  thumbnailUrl: string;
};

export type AlbumInfo = {
  id: number;
  user_id: number;
  name: string;
  description?: string;
  created_at: string;
  updated_at: string;
};

export type CreateAlbumRequest = {
  name: string;
  description?: string;
};

export type UpdateAlbumRequest = {
  name: string;
  description?: string;
};

export type AddImageToAlbumRequest = {
  image_id: string;
};

export type DeletePhotoResponse = {
  message: string;
};

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

export type UploadFormProps = {
  onComplete: () => void;
};
