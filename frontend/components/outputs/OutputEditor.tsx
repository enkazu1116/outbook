// OutputEditor：左テキスト/右プレビューの二画面編集
// - ツールバーで簡易Markdown挿入、メタ情報（種類/公開範囲/タグ）編集
"use client";
import { useState } from 'react';
import { MarkdownRenderer } from './MarkdownRenderer';

export type OutputEditorValue = {
  title: string;
  content: string;
  type: 'summary' | 'note' | 'review' | 'tutorial';
  visibility: 'public' | 'followers' | 'private';
  tags: string[];
};

export function OutputEditor({ initial, onSubmit }: { initial?: Partial<OutputEditorValue>; onSubmit: (value: OutputEditorValue) => void }) {
  const [title, setTitle] = useState(initial?.title ?? '');
  const [content, setContent] = useState(initial?.content ?? '');
  const [type, setType] = useState<OutputEditorValue['type']>(initial?.type ?? 'summary');
  const [visibility, setVisibility] = useState<OutputEditorValue['visibility']>(initial?.visibility ?? 'public');
  const [tags, setTags] = useState((initial?.tags ?? []).join(', '));

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const tagList = tags.split(',').map((t) => t.trim()).filter(Boolean).slice(0, 10);
    onSubmit({ title, content, type, visibility, tags: tagList });
  };

  return (
    <form className="editor" onSubmit={handleSubmit}>
      <div className="editor__toolbar" aria-label="エディタツールバー">
        <button type="button" onClick={() => setContent((c) => c + '\n\n## 見出し')}>H2</button>
        <button type="button" onClick={() => setContent((c) => c + ' **太字** ')}>B</button>
        <button type="button" onClick={() => setContent((c) => c + ' _斜体_ ')}>I</button>
        <button type="button" onClick={() => setContent((c) => c + ' `code` ')}>Code</button>
        <button type="button" onClick={() => setContent((c) => c + '\n- item')}>List</button>
        <button type="button" onClick={() => setContent((c) => c + ' [link](https://example.com) ')}>Link</button>
        <button type="button" onClick={() => setContent((c) => c + '\n![alt](image.jpg)')}>Image</button>
      </div>

      <label className="form__label">タイトル</label>
      <input className="form__input" name="title" value={title} onChange={(e) => setTitle(e.target.value)} required />

      <div className="editor__panes">
        <textarea className="editor__input" name="content" value={content} onChange={(e) => setContent(e.target.value)} />
        <div className="editor__preview md"><MarkdownRenderer markdown={content} /></div>
      </div>

      <div className="editor__meta">
        <div className="field">
          <label className="form__label">種類</label>
          <select className="form__select" value={type} onChange={(e) => setType(e.target.value as any)}>
            <option value="summary">要約</option>
            <option value="note">ノート</option>
            <option value="review">レビュー</option>
            <option value="tutorial">チュートリアル</option>
          </select>
        </div>
        <div className="field">
          <label className="form__label">公開範囲</label>
          <select className="form__select" value={visibility} onChange={(e) => setVisibility(e.target.value as any)}>
            <option value="public">公開</option>
            <option value="followers">フォロワー</option>
            <option value="private">非公開</option>
          </select>
        </div>
        <div className="field">
          <label className="form__label">タグ</label>
          <input className="form__input" value={tags} onChange={(e) => setTags(e.target.value)} placeholder="カンマ区切り 例: architecture, pattern" />
          <p className="form__help">最大10個まで</p>
        </div>
      </div>

      <div className="editor__actions">
        <button className="button button--secondary" type="button">下書き保存</button>
        <button className="button button--primary" type="submit">投稿</button>
      </div>
    </form>
  );
}


