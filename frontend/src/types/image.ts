export interface ContainerImage {
  id: string;
  name: string;
  description: string;
  author: string;
  registry: string;
  namespace: string;
  repository: string;
  tag: string;
  digest: string;
  size: number;
  readme_path: string;
  stars: number;
  visibility: 'public' | 'private';
  platform: string;
  labels: Label[];
  created_at: string;
  updated_at: string;
}

export interface Label {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface ImageListResponse {
  data: ContainerImage[];
  total: number;
}

export interface ImageListRequest {
  page?: number;
  page_size?: number;
  search?: string;
} 