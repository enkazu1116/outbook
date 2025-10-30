// Redux store: UI（テーマ/サイドバー）と ランキングを集中管理
// - ui: 画面全体の表示状態（テーマ/サイドバー）
// - ranking: Amazon/Gihyoの売上ランキング取得状態（source/items/loading/error）
// - fetchRanking: APIからランキングを取得する thunk
import { configureStore, createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';

type ThemeMode = 'light' | 'dark';

type UiState = {
  theme: ThemeMode;
  sidebarCollapsed: boolean;
};

const initialUiState: UiState = {
  theme: 'light',
  sidebarCollapsed: false,
};

const uiSlice = createSlice({
  name: 'ui',
  initialState: initialUiState,
  reducers: {
    toggleSidebar(state) { state.sidebarCollapsed = !state.sidebarCollapsed; },
    setSidebar(state, action: PayloadAction<boolean>) { state.sidebarCollapsed = action.payload; },
    setTheme(state, action: PayloadAction<ThemeMode>) { state.theme = action.payload; },
    toggleTheme(state) { state.theme = state.theme === 'dark' ? 'light' : 'dark'; },
  },
});

type RankingItem = { rank: number; title: string; author?: string; cover?: string; url?: string };
type Source = 'amazon' | 'gihyo';
type RankingState = {
  source: Source;
  items: RankingItem[];
  loading: boolean;
  error?: string | null;
};

const initialRankingState: RankingState = {
  source: 'amazon',
  items: [],
  loading: false,
  error: null,
};

export const fetchRanking = createAsyncThunk<RankingItem[], Source>(
  'ranking/fetch',
  async (source: Source) => {
    const res = await fetch(`/api/ranking?source=${source}`);
    const data = await res.json();
    return data.items ?? [];
  }
);

const rankingSlice = createSlice({
  name: 'ranking',
  initialState: initialRankingState,
  reducers: {
    setSource(state, action: PayloadAction<Source>) { state.source = action.payload; },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchRanking.pending, (state) => { state.loading = true; state.error = null; })
      .addCase(fetchRanking.fulfilled, (state, action) => { state.loading = false; state.items = action.payload; })
      .addCase(fetchRanking.rejected, (state, action) => { state.loading = false; state.error = action.error.message ?? '取得に失敗しました'; });
  },
});

export const { toggleSidebar, setSidebar, setTheme, toggleTheme } = uiSlice.actions;
export const { setSource } = rankingSlice.actions;

export const store = configureStore({
  reducer: {
    ui: uiSlice.reducer,
    ranking: rankingSlice.reducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;


