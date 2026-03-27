import { create } from "zustand";

export const gameStore = create((set, get) => ({
  socket: null,
  messages: [],
  players: [],
  ingamePlayers: [],
  gameState: {
    role: "",
    currentTurn: "",
    topCard: { color: "", value: "" },
    cardInHand: [],
  },

  setSocket: (ws) => set({ socket: ws }),

  sendMessage: (data) => {
    const socket = get().socket;
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(data));
    }
  },

  addMessage: (msg) =>
    set((state) => ({
      messages: [...state.messages, msg],
    })),

  setPlayers: (data) =>
    set(() => ({
      players: [...data],
    })),

  setInGamePlayers: (data) =>
    set((state) => ({
      ingamePlayers: [...state.ingamePlayers, data],
    })),

  resetInGamePlayers: (playerId) =>
    set((state) => ({
      ingamePlayers: state.ingamePlayers.filter((p) => p.playerId !== playerId),
    })),

  setGameStateCards: (cards) =>
    set((state) => ({
      gameState: {
        ...state.gameState,
        cardInHand: cards,
      },
    })),
  setGameStateRole: (role) =>
    set((state) => ({
      gameState: {
        ...state.gameState,
        role: role,
      },
    })),
  setGameStateTop: (top) =>
    set((state) => ({
      gameState: {
        ...state.gameState,
        topCard: top,
      },
    })),
  setGameStateTurn: (turn) =>
    set((state) => ({
      gameState: {
        ...state.gameState,
        currentTurn: turn,
      },
    })),
}));
