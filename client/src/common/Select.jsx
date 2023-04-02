import React from "react";
import "tailwindcss/tailwind.css";

function Select({ listCommands, selectedCommand, onCommandChange }) {
  const handleCommandChange = (e) => {
    onCommandChange(e.target.value);
  };

  return (
    <div className="flex flex-col items-center justify-center">
      <h1 className="mb-4 text-3xl font-bold">Select Command:</h1>
      <select
        className="w-80 h-10 px-3 mb-4 text-base text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
        value={selectedCommand}
        onChange={handleCommandChange}
      >
        {listCommands.map((cmd) => (
          <option key={cmd} value={cmd}>
            {cmd}
          </option>
        ))}
      </select>
      <p>You selected: {selectedCommand}</p>
    </div>
  );
}

export default Select;
