export const API_URL = import.meta.env.VITE_MAKTABA_API_URL || "http://localhost:8080";

export type User = {
  id: string;
  email?: string;
  username?: string;
  firstName?: string;
  lastName?: string;
};

export type Source = {
  id: string;
  title: string;
  subtitle?: string;
  type: string;
  description?: string;
  publisher?: string;
  isbn?: string;
  tags?: string[];
  created_at: string;
};

export type Book = {
  source: Source;
  metadata?: {
    isbn_10?: string;
    isbn_13?: string;
    publisher?: string;
    page_count?: number;
    language?: string;
    cover_url?: string;
  };
  contributors?: Array<{
    id: string;
    name: string;
    role: string;
    position: number;
  }>;
};

export type LibraryItem = {
  id: string;
  user_id: string;
  source_id: string;
  status: string;
  progress_value?: number;
  progress_unit?: string;
  visibility: string;
  created_at: string;
  updated_at: string;
};

export type LibraryItemWithSource = LibraryItem & {
  source: Source;
};

export type Note = {
  id: string;
  user_id: string;
  source_id?: string;
  content: string;
  content_type: string;
  is_public: boolean;
  tags?: string[];
  created_at: string;
};

export type Review = {
  id: string;
  user_id: string;
  source_id: string;
  rating: number;
  content?: string;
  is_public: boolean;
  created_at: string;
};

export type Collection = {
  id: string;
  user_id: string;
  name: string;
  description?: string;
  is_public: boolean;
  source_ids?: string[];
  created_at: string;
};

export type Profile = {
  id: string;
  user_id: string;
  username?: string;
  display_name?: string;
  bio?: string;
  public_profile: boolean;
  created_at: string;
  updated_at: string;
};

type RequestOptions = RequestInit & {
  accessToken?: string | null;
};

export async function apiRequest<T>(path: string, options: RequestOptions = {}) {
  const headers = new Headers(options.headers);
  if (options.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  if (options.accessToken) {
    headers.set("Authorization", `Bearer ${options.accessToken}`);
  }

  const response = await fetch(`${API_URL}${path}`, {
    ...options,
    headers,
    credentials: "include",
  });

  if (!response.ok) {
    const fallback = `Request failed with ${response.status}`;
    let message = fallback;
    try {
      const data = await response.json();
      message = data.error || fallback;
    } catch {
      // Keep fallback for non-JSON errors.
    }
    throw new Error(message);
  }

  if (response.status === 204) {
    return undefined as T;
  }
  return response.json() as Promise<T>;
}

export function getMe(accessToken: string) {
  return apiRequest<User>("/api/me", { accessToken });
}

export function listSources() {
  return apiRequest<Source[]>("/sources?type=book&limit=100");
}

export function createBook(accessToken: string, payload: unknown) {
  return apiRequest<Book>("/api/sources/books", {
    method: "POST",
    accessToken,
    body: JSON.stringify(payload),
  });
}

export function listLibrary(accessToken: string) {
  return apiRequest<LibraryItemWithSource[]>("/api/library/items/with-sources?limit=50", {
    accessToken,
  });
}

export function addLibraryItem(accessToken: string, sourceID: string) {
  return apiRequest<LibraryItem>("/api/library/items", {
    method: "POST",
    accessToken,
    body: JSON.stringify({
      source_id: sourceID,
      status: "to_consume",
      visibility: "private",
    }),
  });
}

export function updateLibraryItem(accessToken: string, itemID: string, payload: unknown) {
  return apiRequest<LibraryItem>(`/api/library/items/${encodeURIComponent(itemID)}`, {
    method: "PUT",
    accessToken,
    body: JSON.stringify(payload),
  });
}

export function deleteLibraryItem(accessToken: string, itemID: string) {
  return apiRequest<void>(`/api/library/items/${encodeURIComponent(itemID)}`, {
    method: "DELETE",
    accessToken,
  });
}

export function listNotes(accessToken: string) {
  return apiRequest<Note[]>("/api/notes?limit=50", { accessToken });
}

export function listReviews(accessToken: string) {
  return apiRequest<Review[]>("/api/reviews?limit=50", { accessToken });
}

export function listCollections(accessToken: string) {
  return apiRequest<Collection[]>("/api/collections?limit=50", { accessToken });
}

export function createNote(accessToken: string, payload: unknown) {
  return apiRequest<Note>("/api/notes", {
    method: "POST",
    accessToken,
    body: JSON.stringify(payload),
  });
}

export function deleteNote(accessToken: string, noteID: string) {
  return apiRequest<void>(`/api/notes/${encodeURIComponent(noteID)}`, {
    method: "DELETE",
    accessToken,
  });
}

export function createReview(accessToken: string, payload: unknown) {
  return apiRequest<Review>("/api/reviews", {
    method: "POST",
    accessToken,
    body: JSON.stringify(payload),
  });
}

export function deleteReview(accessToken: string, reviewID: string) {
  return apiRequest<void>(`/api/reviews/${encodeURIComponent(reviewID)}`, {
    method: "DELETE",
    accessToken,
  });
}

export function createCollection(accessToken: string, payload: unknown) {
  return apiRequest<Collection>("/api/collections", {
    method: "POST",
    accessToken,
    body: JSON.stringify(payload),
  });
}

export function deleteCollection(accessToken: string, collectionID: string) {
  return apiRequest<void>(`/api/collections/${encodeURIComponent(collectionID)}`, {
    method: "DELETE",
    accessToken,
  });
}

export function getProfile(accessToken: string) {
  return apiRequest<Profile>("/api/profile", { accessToken });
}

export function updateProfile(accessToken: string, payload: unknown) {
  return apiRequest<Profile>("/api/profile", {
    method: "PUT",
    accessToken,
    body: JSON.stringify(payload),
  });
}

export function getPublicProfile(username: string) {
  return apiRequest<Profile>(`/users/${encodeURIComponent(username)}/profile`);
}

export function listPublicLibrary(username: string) {
  return apiRequest<LibraryItemWithSource[]>(
    `/users/${encodeURIComponent(username)}/library/with-sources?limit=50`
  );
}

export function getBook(sourceID: string) {
  return apiRequest<Book>(`/sources/books/${encodeURIComponent(sourceID)}`);
}

export function listPublicNotesByUser(userID: string) {
  return apiRequest<Note[]>(`/notes?user_id=${encodeURIComponent(userID)}&limit=50`);
}

export function listPublicReviewsByUser(userID: string) {
  return apiRequest<Review[]>(`/reviews?user_id=${encodeURIComponent(userID)}&limit=50`);
}

export function listPublicCollectionsByUser(userID: string) {
  return apiRequest<Collection[]>(`/collections?user_id=${encodeURIComponent(userID)}&limit=50`);
}
