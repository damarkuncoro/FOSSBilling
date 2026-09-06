import React from 'react';
import { BookOpen, Search, ArrowLeft, ThumbsUp, Eye, ChevronRight } from 'lucide-react';
import { useKnowledgebase } from '../hooks/useKnowledgebase';

export const Knowledgebase: React.FC = () => {
  const {
    articles,
    categories,
    search,
    setSearch,
    selectedCategory,
    setSelectedCategory,
    activeArticle,
    setActiveArticle,
  } = useKnowledgebase();

  if (activeArticle) {
    return (
      <div className="max-w-4xl mx-auto px-4 py-8">
        <button
          onClick={() => setActiveArticle(null)}
          className="inline-flex items-center gap-2 text-sm font-medium text-gray-500 hover:text-indigo-600 mb-6 transition-colors"
        >
          <ArrowLeft className="w-4 h-4" /> Back to Knowledgebase
        </button>

        <article className="bg-white rounded-2xl border border-gray-100 shadow-sm p-8 space-y-6">
          <div className="space-y-2">
            <span className="px-2.5 py-1 rounded-full text-xs font-semibold bg-indigo-50 text-indigo-700">
              {activeArticle.category}
            </span>
            <h1 className="text-2xl font-bold text-gray-900 tracking-tight">{activeArticle.title}</h1>
            <div className="flex items-center gap-4 text-xs text-gray-400 pt-1 border-b border-gray-100 pb-4">
              <span className="flex items-center gap-1"><Eye className="w-3.5 h-3.5" /> {activeArticle.views} views</span>
              <span>Updated: {new Date(activeArticle.updated_at).toLocaleDateString()}</span>
            </div>
          </div>

          <div className="prose prose-indigo max-w-none text-gray-700 leading-relaxed text-sm">
            <p className="font-medium text-gray-800 text-base">{activeArticle.summary}</p>
            <div className="mt-4 p-4 bg-gray-50 rounded-xl border border-gray-100 font-mono text-xs">
              {activeArticle.content}
            </div>
          </div>

          <div className="pt-6 border-t border-gray-100 flex items-center justify-between">
            <span className="text-xs text-gray-500">Was this article helpful?</span>
            <button className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-emerald-50 text-emerald-700 hover:bg-emerald-100 rounded-lg text-xs font-semibold transition-colors">
              <ThumbsUp className="w-3.5 h-3.5" /> Helpful ({activeArticle.helpful_count})
            </button>
          </div>
        </article>
      </div>
    );
  }

  return (
    <div className="max-w-5xl mx-auto px-4 py-8 space-y-8">
      <div className="text-center max-w-2xl mx-auto space-y-3">
        <h1 className="text-3xl font-extrabold text-gray-900 tracking-tight flex items-center justify-center gap-2.5">
          <BookOpen className="w-8 h-8 text-indigo-600" /> Knowledgebase & Tutorials
        </h1>
        <p className="text-sm text-gray-500">
          Find answers, server setup guides, and billing walkthroughs.
        </p>

        <div className="relative max-w-lg mx-auto pt-2">
          <Search className="w-5 h-5 absolute left-4 top-1/2 translate-y-1 text-gray-400" />
          <input
            type="text"
            placeholder="Search guides (e.g. DNS, SSL, cPanel)..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full pl-12 pr-4 py-3 bg-white border border-gray-200 rounded-2xl shadow-sm text-sm focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none"
          />
        </div>
      </div>

      <div className="flex flex-wrap items-center justify-center gap-2">
        {categories.map((cat) => (
          <button
            key={cat}
            onClick={() => setSelectedCategory(cat)}
            className={`px-4 py-2 rounded-xl text-xs font-semibold capitalize transition-all ${
              selectedCategory === cat
                ? 'bg-indigo-600 text-white shadow-sm'
                : 'bg-white border border-gray-200 text-gray-600 hover:bg-gray-50'
            }`}
          >
            {cat}
          </button>
        ))}
      </div>

      {articles.length === 0 ? (
        <div className="bg-white rounded-2xl border border-dashed border-gray-200 p-12 text-center space-y-3">
          <BookOpen className="w-10 h-10 text-gray-400 mx-auto" />
          <h3 className="text-base font-bold text-gray-900">No Articles Found</h3>
          <p className="text-xs text-gray-500 max-w-sm mx-auto">
            No knowledgebase tutorials match &quot;{search}&quot;. Try searching with different keywords or clearing filters.
          </p>
          <button
            onClick={() => { setSearch(''); setSelectedCategory('all'); }}
            className="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 text-xs font-semibold rounded-xl transition-colors"
          >
            Clear Search & Filters
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {articles.map((art) => (
            <div
              key={art.id}
              onClick={() => setActiveArticle(art)}
              className="group bg-white p-6 rounded-2xl border border-gray-100 shadow-sm hover:border-indigo-200 hover:shadow-md cursor-pointer transition-all flex flex-col justify-between"
            >
              <div>
                <span className="text-[11px] font-semibold text-indigo-600 uppercase tracking-wider">
                  {art.category}
                </span>
                <h3 className="text-base font-bold text-gray-900 group-hover:text-indigo-600 mt-1 transition-colors">
                  {art.title}
                </h3>
                <p className="text-xs text-gray-500 mt-2 line-clamp-2">{art.summary}</p>
              </div>

              <div className="flex items-center justify-between pt-4 mt-4 border-t border-gray-100 text-xs text-gray-400">
                <span>{art.views} reads</span>
                <span className="text-indigo-600 font-semibold flex items-center gap-1 group-hover:translate-x-1 transition-transform">
                  Read Article <ChevronRight className="w-3.5 h-3.5" />
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
