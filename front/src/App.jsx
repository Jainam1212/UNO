import { useEffect } from "react";
import { gameStore } from "./store/sockerStore";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Home } from "./components/Home";
import { useCookies } from "react-cookie";

function App() {
  const setSocket = gameStore((s) => s.setSocket);
  const addMessage = gameStore((s) => s.addMessage);
  const addPlayer = gameStore((s) => s.setPlayers);
  const addInGamePlayer = gameStore((s) => s.setInGamePlayers);
  const [, setCookie] = useCookies(["userinfo-cookie"]);

  useEffect(() => {
    const ws = new WebSocket("ws://192.168.1.23:8080/playerHandler");

    ws.onopen = () => {
      console.log("Connected");
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      if (data.type === "user_info" && data.playerId) {
        setCookie("userinfo-cookie", data.playerId);
      }
      switch (data.type) {
        case "message":
          addMessage((prev) => [...prev, data]);
          break;
        case "ingame_player_list":
          addInGamePlayer(data.inGamePlayerIds);
          break;
        case "players_chat":
          addMessage({
            playerName: data.playerName,
            playerMessage: data.playerMessage,
          });
          break;
        case "add_player_list":
          addPlayer(data.newPlayerIds);
          break;

        default:
          break;
      }
    };

    // eslint-disable-next-line react-hooks/set-state-in-effect
    setSocket(ws);

    return () => ws.close();
  }, []);

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
