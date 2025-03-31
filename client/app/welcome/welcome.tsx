import { useState } from "react";

export function Welcome() {
  const fetchData = () => {
    fetch(`http://localhost:${import.meta.env.VITE_PORT}/`)
      .then((response) => response.text())
      .then((data) => setMessage(data))
      .catch((error) => console.error("Error fetching data:", error));
  };
  const [message, setMessage] = useState<string>("");
  return (
    <div className="flex flex-col gap-2">
      <button onClick={fetchData}>Fetch Data</button>
      <p>{message}</p>
    </div>
  );
}
