import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import type { FormEvent } from "react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import {
  AddBookCard,
  CollectionsCard,
  type DashboardData,
  DashboardHeader,
  DashboardIntro,
  DashboardStats,
  LibraryPanel,
  RecentNotesCard,
  RecentReviewsCard,
  SourceActivityCard,
  SourcesPanel,
} from "~/components/dashboard/dashboard-panels";
import {
  addLibraryItem,
  createBook,
  createCollection,
  createNote,
  createReview,
  deleteCollection,
  deleteLibraryItem,
  deleteNote,
  deleteReview,
  listCollections,
  listLibrary,
  listNotes,
  listReviews,
  listSources,
  updateLibraryItem,
} from "~/lib/api";
import { useAuthStore } from "~/lib/auth";

const emptyData: DashboardData = {
  sources: [],
  library: [],
  notes: [],
  reviews: [],
  collections: [],
};

export function normalizeDashboardData(data: Partial<DashboardData>): DashboardData {
  return {
    sources: Array.isArray(data.sources) ? data.sources : [],
    library: Array.isArray(data.library) ? data.library : [],
    notes: Array.isArray(data.notes) ? data.notes : [],
    reviews: Array.isArray(data.reviews) ? data.reviews : [],
    collections: Array.isArray(data.collections) ? data.collections : [],
  };
}

export default function DashboardPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { isAuthenticated, isLoading, user, accessToken, logout } = useAuthStore();
  const [error, setError] = useState<string | null>(null);
  const [bookTitle, setBookTitle] = useState("");
  const [bookAuthor, setBookAuthor] = useState("");
  const [bookISBN, setBookISBN] = useState("");
  const [selectedSourceID, setSelectedSourceID] = useState("");
  const [noteContent, setNoteContent] = useState("");
  const [notePublic, setNotePublic] = useState(false);
  const [reviewRating, setReviewRating] = useState("5");
  const [reviewContent, setReviewContent] = useState("");
  const [collectionName, setCollectionName] = useState("");
  const [collectionPublic, setCollectionPublic] = useState(false);

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate("/login");
    }
  }, [isAuthenticated, isLoading, navigate]);

  const dashboardQuery = useQuery({
    queryKey: ["dashboard", accessToken],
    enabled: Boolean(isAuthenticated && accessToken),
    queryFn: async () => {
      const token = accessToken as string;
      const [sources, library, notes, reviews, collections] = await Promise.all([
        listSources(),
        listLibrary(token),
        listNotes(token),
        listReviews(token),
        listCollections(token),
      ]);
      return normalizeDashboardData({ sources, library, notes, reviews, collections });
    },
  });

  const data = normalizeDashboardData(dashboardQuery.data || emptyData);

  useEffect(() => {
    if (!selectedSourceID && data.library.length > 0) {
      setSelectedSourceID(data.library[0].source_id);
    }
  }, [data.library, selectedSourceID]);

  const invalidateDashboard = () => queryClient.invalidateQueries({ queryKey: ["dashboard"] });

  const createBookMutation = useMutation({
    mutationFn: async () => {
      if (!accessToken || !bookTitle.trim()) return;
      const book = await createBook(accessToken, {
        title: bookTitle.trim(),
        isbn_13: bookISBN.trim() || undefined,
        contributors: bookAuthor.trim() ? [{ name: bookAuthor.trim(), role: "author" }] : [],
      });
      await addLibraryItem(accessToken, book.source.id);
    },
    onSuccess: async () => {
      setBookTitle("");
      setBookAuthor("");
      setBookISBN("");
      await invalidateDashboard();
    },
    onError: (err) => setError(err instanceof Error ? err.message : "Failed to create book"),
  });

  const activityMutation = useMutation({
    mutationFn: async (action: () => Promise<unknown>) => action(),
    onSuccess: async () => invalidateDashboard(),
    onError: (err) => setError(err instanceof Error ? err.message : "Action failed"),
  });

  const handleLogout = () => {
    logout();
    queryClient.clear();
    navigate("/");
  };

  const handleCreateBook = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError(null);
    createBookMutation.mutate();
  };

  const handleCreateNote = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!accessToken || !selectedSourceID || !noteContent.trim()) return;
    setError(null);
    activityMutation.mutate(async () => {
      await createNote(accessToken, {
        source_id: selectedSourceID,
        content: noteContent.trim(),
        content_type: "note",
        is_public: notePublic,
      });
      setNoteContent("");
      setNotePublic(false);
    });
  };

  const handleCreateReview = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!accessToken || !selectedSourceID) return;
    setError(null);
    activityMutation.mutate(async () => {
      await createReview(accessToken, {
        source_id: selectedSourceID,
        rating: Number(reviewRating),
        content: reviewContent.trim() || undefined,
        is_public: true,
      });
      setReviewRating("5");
      setReviewContent("");
    });
  };

  const handleCreateCollection = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!accessToken || !selectedSourceID || !collectionName.trim()) return;
    setError(null);
    activityMutation.mutate(async () => {
      await createCollection(accessToken, {
        name: collectionName.trim(),
        is_public: collectionPublic,
        source_ids: [selectedSourceID],
      });
      setCollectionName("");
      setCollectionPublic(false);
    });
  };

  const runProtected = (action: () => Promise<unknown>) => {
    if (!accessToken) return;
    setError(null);
    activityMutation.mutate(action);
  };

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-slate-50">
        <Loader2 className="h-8 w-8 animate-spin text-emerald-600" />
      </div>
    );
  }
  if (!isAuthenticated) return null;

  const librarySourceIDs = new Set(data.library.map((item) => item.source_id));
  const librarySources = data.library.map((item) => ({ item, source: item.source }));

  return (
    <div className="min-h-screen bg-slate-50">
      <DashboardHeader
        displayName={user.username || user.email || "reader"}
        onLogout={handleLogout}
      />
      <main className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
        <DashboardIntro />
        {(error || dashboardQuery.error) && (
          <div className="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
            {error ||
              (dashboardQuery.error instanceof Error
                ? dashboardQuery.error.message
                : "Failed to load dashboard")}
          </div>
        )}
        <DashboardStats data={data} />
        <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_380px]">
          <div className="space-y-6">
            <LibraryPanel
              accessToken={accessToken}
              items={data.library}
              loading={dashboardQuery.isLoading}
              onDelete={(id) => runProtected(() => deleteLibraryItem(accessToken as string, id))}
              onUpdate={(id, payload) =>
                runProtected(() => updateLibraryItem(accessToken as string, id, payload))
              }
            />
            <SourcesPanel
              accessToken={accessToken}
              librarySourceIDs={librarySourceIDs}
              loading={dashboardQuery.isLoading}
              onAdd={(id) => runProtected(() => addLibraryItem(accessToken as string, id))}
              sources={data.sources}
            />
          </div>
          <aside className="space-y-6">
            <AddBookCard
              author={bookAuthor}
              isbn={bookISBN}
              saving={createBookMutation.isPending}
              title={bookTitle}
              onAuthorChange={setBookAuthor}
              onISBNChange={setBookISBN}
              onSubmit={handleCreateBook}
              onTitleChange={setBookTitle}
            />
            <SourceActivityCard
              collectionName={collectionName}
              collectionPublic={collectionPublic}
              librarySources={librarySources}
              noteContent={noteContent}
              notePublic={notePublic}
              reviewContent={reviewContent}
              reviewRating={reviewRating}
              saving={activityMutation.isPending}
              selectedSourceID={selectedSourceID}
              onCollectionNameChange={setCollectionName}
              onCollectionPublicChange={setCollectionPublic}
              onNoteContentChange={setNoteContent}
              onNotePublicChange={setNotePublic}
              onReviewContentChange={setReviewContent}
              onReviewRatingChange={setReviewRating}
              onSelectedSourceChange={setSelectedSourceID}
              onSubmitCollection={handleCreateCollection}
              onSubmitNote={handleCreateNote}
              onSubmitReview={handleCreateReview}
            />
            <RecentNotesCard
              notes={data.notes}
              onDelete={(id) => runProtected(() => deleteNote(accessToken as string, id))}
            />
            <RecentReviewsCard
              reviews={data.reviews}
              onDelete={(id) => runProtected(() => deleteReview(accessToken as string, id))}
            />
            <CollectionsCard
              collections={data.collections}
              onDelete={(id) => runProtected(() => deleteCollection(accessToken as string, id))}
            />
          </aside>
        </div>
      </main>
    </div>
  );
}
