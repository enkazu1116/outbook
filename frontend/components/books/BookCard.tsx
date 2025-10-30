"use client";
// 書籍カード：タイトル/著者/タグを表示し、クリック/Enterで詳細に遷移できる
export type BookCardProps = {
  title: string;
  author: string;
  cover?: string;
  tags?: string[];
  onClick?: () => void;
};

import type { KeyboardEvent } from 'react';
import { memo } from 'react';

function BookCardBase({ title, author, cover, tags, onClick }: BookCardProps) {
  return (
    <article className="book-card" role="button" tabIndex={0} onClick={onClick} onKeyDown={(e: KeyboardEvent<HTMLElement>) => { if (e.key === 'Enter') onClick?.(); }}>
      <img className="book-card__cover" src={cover || '/covers/placeholder.jpg'} alt="書籍カバー" />
      <div className="book-card__body">
        <h3 className="book-card__title">{title}</h3>
        <p className="book-card__author">{author}</p>
        <div className="book-card__meta">
          {tags?.map((t) => (<span className="tag" key={t}>{t}</span>))}
        </div>
      </div>
    </article>
  );
}

export const BookCard = memo(BookCardBase);


