// ルートレイアウト
// - 全CSS読み込み、テーマ初期化スクリプト、Redux Provider の設置
import type { Metadata } from 'next';
import '@/styles/base.css';
import '@/styles/layout.css';
import '@/styles/components.css';
import '@/styles/pages.css';
import '@/styles/themes.css';
import StoreProvider from '@/components/providers/StoreProvider';

export const metadata: Metadata = {
  title: 'OutBook',
  description: 'Learning outputs platform',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ja">
      <body className="theme-light">
        <script
          dangerouslySetInnerHTML={{
            __html: `
              (function(){
                try {
                  var s = localStorage.getItem('theme');
                  var m = s ? s : (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
                  document.body.classList.remove('theme-light','theme-dark');
                  document.body.classList.add(m === 'dark' ? 'theme-dark' : 'theme-light');
                } catch (e) {}
              })();
            `,
          }}
        />
        <StoreProvider>
          {children}
        </StoreProvider>
      </body>
    </html>
  );
}


