import { useState } from "react";

import "../styles/home.css";
import { gameStore } from "../store/sockerStore";
import { customFetch } from "../utils/useFetch";
import { useCookies } from "react-cookie";

export const Home = () => {
  const [page, setPage] = useState("home"); // home | room
  const [modal, setModal] = useState(null); // null | create | join
  const [role, setRole] = useState("");
  const [chatInput, setChatInput] = useState("");

  const [playerName, setPlayerName] = useState("");
  const [gameId, setGameId] = useState("");
  const [maxPlayers, setMaxPlayers] = useState(6);

  const [joinPlayerName, setJoinPlayerName] = useState("");
  const [joinGameId, setJoinGameId] = useState("");

  const playerList = gameStore((s) => s.players);
  const ingamePlayerList = gameStore((s) => s.ingamePlayers);
  const socketMessenger = gameStore((s) => s.sendMessage);
  const chat = gameStore((s) => s.messages);

  const [cookies] = useCookies(["userinfo-cookie"]);

  const enterRoom = () => {
    setModal(null);
    setPage("room");
    const pid = cookies["userinfo-cookie"];
    if (!pid) {
      return;
    }
    customFetch("http://192.168.1.6:8080/goapi/v1/createGame", "post", {
      playerId: Number(pid),
      gameId: Number(gameId),
      playerName,
      maxPlayers: Number(maxPlayers),
    });
    setRole("owner");
  };
  const joinRoom = () => {
    setModal(null);
    setPage("room");
    const pid = cookies["userinfo-cookie"];
    if (!pid) {
      return;
    }
    customFetch("http://192.168.1.6:8080/goapi/v1/joinGame", "post", {
      playerId: Number(pid),
      gameId: Number(joinGameId),
      playerName: joinPlayerName,
    });
    setRole("player");
  };

  const leaveGame = () => {
    setModal(null);
    setPage("home");
    const pid = cookies["userinfo-cookie"];
    if (!pid) {
      return;
    }
    customFetch("http://192.168.1.6:8080/goapi/v1/leaveGame", "post", {
      playerId: Number(pid),
      gameId: Number(joinGameId),
    });
  };

  const sendChat = () => {
    if (!chatInput.trim()) return;
    socketMessenger({
      messageType: "message",
      userName: playerName || joinPlayerName + " said",
      userMessage: chatInput.trim(),
    });
    setChatInput("");
  };

  return (
    <>
      <div className="app">
        <div className="blob blob-1" />
        <div className="blob blob-2" />
        <div className="blob blob-3" />

        <div className="content">
          {/* NAV */}
          <nav>
            <div className="logo">UNO</div>
            <div className="nav-right">
              <div className="badge">
                <span>●</span>
                {playerList.length} online
              </div>
              {page === "room" && (
                <div className="badge">
                  <span>🃏</span>Room: {gameId}
                </div>
              )}
            </div>
          </nav>

          {/* HOME */}
          {page === "home" && (
            <div className="home">
              <div className="hero">
                <div className="hero-eyebrow">The Classic Card Game</div>
                <div className="hero-title">UNO</div>
                <div className="hero-sub">
                  Challenge your friends. Draw four. Skip turns. Win glory.
                </div>
              </div>

              <div className="cards">
                {/* Create Room */}
                <div
                  className="card card-create"
                  onClick={() => {
                    setModal("create");
                  }}
                >
                  <div className="card-icon">🎮</div>
                  <h3>Create Room</h3>
                  <p>
                    Start a new game and invite your friends with a unique room
                    code.
                  </p>
                  <button className="btn btn-red">Create Room</button>
                </div>

                {/* Join Room */}
                <div
                  className="card card-join"
                  onClick={() => setModal("join")}
                >
                  <div className="card-icon">🔗</div>
                  <h3>Join Room</h3>
                  <p>
                    Enter a room code to jump into an existing game with
                    friends.
                  </p>
                  <button className="btn btn-blue">Join Room</button>
                </div>
              </div>
            </div>
          )}

          {/* ROOM */}
          {page === "room" && (
            <div className="room-layout">
              {/* Left sidebar – players */}
              <div className="sidebar">
                <div className="sidebar-title">
                  Players ({ingamePlayerList.length}/6)
                </div>
                <div className="player-list">
                  {ingamePlayerList.map((p, i) => (
                    <div
                      key={i}
                      className={`player-item ${p[0].playerRole === "host" ? "host" : p[0].playerRole === "ready" ? "active" : ""}`}
                    >
                      <div className="avatar">👤</div>
                      <div className="player-info">
                        <div className="player-name">{p[0].playerName}</div>{" "}
                        <div className="player-tag">{p[0].playerRole}</div>
                      </div>
                      {p[0].playerRole === "host" ? (
                        <span className="host-crown">👑</span>
                      ) : (
                        <div className="online-dot" />
                      )}
                    </div>
                  ))}
                </div>
              </div>

              {/* Center */}
              <div className="room-center">
                <div className="room-header">
                  <div className="room-code-label">Room Code</div>
                  <div className="room-code">{gameId}</div>
                </div>

                <div className="card-pile">
                  <div className="uno-card uno-card-1">7</div>
                  <div className="uno-card uno-card-2">+2</div>
                  <div className="uno-card uno-card-3">↺</div>
                </div>

                <div className="status-bar">
                  <div className="status-dot" />
                  Waiting for players… ({ingamePlayerList.length}/6 joined)
                </div>

                <div className="room-actions">
                  <button
                    className="btn-start"
                    disabled={role !== "owner"}
                    style={{
                      cursor: role !== "owner" ? "not-allowed" : "pointer",
                    }}
                  >
                    Start Game
                  </button>
                  <button
                    className="btn-leave"
                    onClick={() => {
                      leaveGame();
                    }}
                  >
                    Leave Room
                  </button>
                </div>
              </div>

              {/* Right sidebar – chat */}
              <div className="sidebar sidebar-right">
                <div className="sidebar-title">Room Chat</div>
                <div className="chat-section">
                  <div className="chat-messages">
                    {chat.map((m, i) => (
                      <div key={i} className="chat-msg">
                        <div className={`chat-msg-name`}>{m.playerName}</div>
                        <div className="chat-msg-text">{m.playerMessage}</div>
                      </div>
                    ))}
                  </div>
                  <div className="chat-input-row">
                    <input
                      placeholder="Say something..."
                      value={chatInput}
                      onChange={(e) => setChatInput(e.target.value)}
                      onKeyDown={(e) => e.key === "Enter" && sendChat()}
                    />
                    <button className="chat-send" onClick={sendChat}>
                      ↑
                    </button>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* MODALS */}
        {modal === "create" && (
          <div className="overlay" onClick={() => setModal(null)}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
              <button className="modal-close" onClick={() => setModal(null)}>
                ✕
              </button>
              <h2 style={{ color: "var(--red)" }}>Create Room</h2>
              <div className="modal-sub">
                Set up your game and share the code
              </div>
              <div className="field">
                <label>Your Name</label>
                <input
                  placeholder="Enter your name…"
                  value={playerName}
                  onChange={(e) => setPlayerName(e.target.value)}
                />
              </div>
              <div className="field">
                <label>Room ID</label>
                <input
                  placeholder="Enter an id for your room"
                  value={gameId}
                  type="number"
                  onChange={(e) => setGameId(e.target.value)}
                />
              </div>
              <div className="field">
                <label>Max Players</label>
                <input
                  placeholder="2 – 6 players"
                  value={maxPlayers}
                  onChange={(e) => {
                    setMaxPlayers(e.target.value);
                  }}
                />
              </div>
              <button className="btn btn-red" onClick={enterRoom}>
                Create &amp; Enter Room
              </button>
            </div>
          </div>
        )}

        {modal === "join" && (
          <div className="overlay" onClick={() => setModal(null)}>
            <div className="modal" onClick={(e) => e.stopPropagation()}>
              <button className="modal-close" onClick={() => setModal(null)}>
                ✕
              </button>
              <h2 style={{ color: "var(--blue)" }}>Join Room</h2>
              <div className="modal-sub">
                Enter the room code from your friend
              </div>
              <div className="field">
                <label>Your Name</label>
                <input
                  placeholder="Enter your name…"
                  value={joinPlayerName}
                  onChange={(e) => setJoinPlayerName(e.target.value)}
                />
              </div>
              <div className="field">
                <label>Room Code</label>
                <input
                  placeholder="e.g. UNO42"
                  value={joinGameId}
                  onChange={(e) => setJoinGameId(e.target.value.toUpperCase())}
                  style={{
                    letterSpacing: "4px",
                    fontFamily: "'Bebas Neue', sans-serif",
                    fontSize: "20px",
                  }}
                />
              </div>
              <button className="btn btn-blue" onClick={joinRoom}>
                Join Room
              </button>
            </div>
          </div>
        )}
      </div>
    </>
  );
};
