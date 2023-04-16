import TodoList from "./pages/TodoList";
import Wallet from "./pages/Wallet";
import { Routes, Route } from "react-router-dom";

const App = () => {
  return (
    <Routes>
      <Route path="/">
        <Route index element={<Wallet />} />
        {/* <Route path="user/:userName" element={<ResultCommit />} /> */}
      </Route>
      <Route path="/todo-app">
        <Route index element={<TodoList />} />
      </Route>
    </Routes>
  );
};

export default App;
