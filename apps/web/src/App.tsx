import { useEffect } from 'react';
import { useArchiveStore } from './store/useArchiveStore';

export default function App() {
  const { seriesList, fetchSeriesList, isAdmin, logoutAdmin } = useArchiveStore();

  useEffect(() => {
    fetchSeriesList();
  }, [fetchSeriesList]);

  return (
    <div className="w-full h-screen bg-velvet-bg flex flex-col items-center justify-center p-6 select-none">
      <div className="text-center space-y-2 mb-8">
        <h1 className="text-velvet-primary text-3xl font-bold tracking-widest animate-pulse">
          THE VELVET ARCHIVE
        </h1>
        <p className="text-velvet-secondary text-sm font-mono">
          Data Interface Layer Operational • Status: {isAdmin ? "ADMIN KEY ENGAGED" : "PUBLIC OBSERVER"}
        </p>
      </div>

      <div className="w-full max-w-md border border-velvet-secondary/30 rounded p-4 bg-black/40 backdrop-blur-md">
        <h2 className="text-velvet-text text-sm font-bold uppercase tracking-wider mb-3 border-b border-velvet-secondary/20 pb-1">
          Detected Bookshelves ({seriesList.length})
        </h2>

        {seriesList.length === 0 ? (
          <p className="text-xs font-mono text-velvet-secondary italic">
            No series records loaded. Populate your Postgres database via admin routing panel or seed script.
          </p>
        ) : (
          <ul className="space-y-2 max-h-48 overflow-y-auto pr-2">
            {seriesList.map((series) => (
              <li 
                key={series.id} 
                className="text-xs font-mono p-2 border border-velvet-secondary/10 bg-velvet-bg/60 rounded flex justify-between items-center"
              >
                <span className="text-velvet-text font-bold">{series.title}</span>
                <span className="text-velvet-secondary italic text-[10px]">by {series.author}</span>
              </li>
            ))}
          </ul>
        )}

        {isAdmin && (
          <button 
            onClick={logoutAdmin}
            className="mt-4 w-full py-1 text-xs font-mono tracking-wider border border-velvet-danger text-velvet-danger hover:bg-velvet-danger/20 transition-all rounded uppercase"
          >
            Drop Admin Key
          </button>
        )}
      </div>
    </div>
  );
}