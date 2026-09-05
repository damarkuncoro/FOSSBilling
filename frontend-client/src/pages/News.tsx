import React from 'react';
import { Newspaper, Calendar, User, ArrowLeft, ChevronRight, AlertTriangle, Sparkles } from 'lucide-react';
import { useNewsArticles } from '../hooks/useNewsArticles';

export const News: React.FC = () => {
  const { articles, selectedArticle, setSelectedArticle } = useNewsArticles();

  if (selectedArticle) {
    return (
      <div className="max-w-3xl mx-auto px-4 py-8">
        <button
          onClick={() => setSelectedArticle(null)}
          className="inline-flex items-center gap-2 text-sm font-medium text-gray-500 hover:text-indigo-600 mb-6 transition-colors"
        >
          <ArrowLeft className="w-4 h-4" /> Back to Announcements
        </button>

        <article className="bg-white rounded-2xl border border-gray-100 shadow-sm p-8 space-y-6">
          <div className="space-y-3 pb-4 border-b border-gray-100">
            <span className={`px-2.5 py-1 rounded-full text-xs font-semibold capitalize ${
              selectedArticle.category === 'maintenance' ? 'bg-amber-50 text-amber-700' : 'bg-indigo-50 text-indigo-700'
            }`}>
              {selectedArticle.category}
            </span>
            <h1 className="text-2xl font-bold text-gray-900 tracking-tight">{selectedArticle.title}</h1>
            <div className="flex items-center gap-4 text-xs text-gray-400">
              <span className="flex items-center gap-1.5"><Calendar className="w-3.5 h-3.5" /> {new Date(selectedArticle.published_at).toLocaleDateString()}</span>
              <span className="flex items-center gap-1.5"><User className="w-3.5 h-3.5" /> {selectedArticle.author}</span>
            </div>
          </div>

          <div className="text-sm text-gray-700 leading-relaxed space-y-4">
            <p>{selectedArticle.content}</p>
          </div>
        </article>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8 space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 tracking-tight flex items-center gap-2.5">
          <Newspaper className="w-7 h-7 text-indigo-600" /> Announcements & Network Status
        </h1>
        <p className="text-sm text-gray-500 mt-1">
          Stay updated with platform releases, maintenance schedules, and security advisories.
        </p>
      </div>

      <div className="space-y-4">
        {articles.map((art) => (
          <div
            key={art.id}
            onClick={() => setSelectedArticle(art)}
            className="group bg-white p-6 rounded-2xl border border-gray-100 shadow-sm hover:border-indigo-200 hover:shadow-md cursor-pointer transition-all flex items-start justify-between gap-4"
          >
            <div className="space-y-2">
              <div className="flex items-center gap-2">
                <span className={`px-2 py-0.5 rounded-full text-[11px] font-semibold capitalize flex items-center gap-1 ${
                  art.category === 'maintenance' ? 'bg-amber-50 text-amber-700' : 'bg-indigo-50 text-indigo-700'
                }`}>
                  {art.category === 'maintenance' ? <AlertTriangle className="w-3 h-3" /> : <Sparkles className="w-3 h-3" />}
                  {art.category}
                </span>
                <span className="text-xs text-gray-400">
                  {new Date(art.published_at).toLocaleDateString()}
                </span>
              </div>
              <h3 className="text-base font-bold text-gray-900 group-hover:text-indigo-600 transition-colors">
                {art.title}
              </h3>
              <p className="text-xs text-gray-500 line-clamp-2">{art.content}</p>
            </div>

            <ChevronRight className="w-5 h-5 text-gray-400 group-hover:text-indigo-600 group-hover:translate-x-1 transition-all shrink-0 my-auto" />
          </div>
        ))}
      </div>
    </div>
  );
};
