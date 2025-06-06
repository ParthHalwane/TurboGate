import { useEffect, useState } from "react";
import RouteForm from "./components/RouteForm";
import RouteList, { type Route } from "./components/RouteList";
import { Toaster } from "react-hot-toast";

function App() {
  const [status, setStatus] = useState<null | string>(null);
  const [routes, setRoutes] = useState<Route[]>([]);

  const fetchRoutes = async () => {
    try {
      const res = await fetch(
        `${import.meta.env.VITE_TURBOGATE_BACKEND_URL}/api/routes`
      );
      const data = await res.json();
      console.log("Fetched route list:", data);
      setRoutes(data);
    } catch (err: unknown) {
      setStatus("❌ Failed to fetch routes");
      console.error("Error fetching routes:", err,status);
    }
  };

  useEffect(() => {
    fetchRoutes();
  });
  
  return (
    <>
      <Toaster position="top-right" reverseOrder={false} />
      <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-900 to-slate-900 text-white p-8">
        <div className="bg-slate-80 bg-opacity-60 backdrop-blur-md rounded-2xl shadow-xl shadow-pink-500/40 p-8 w-full max-w-2xl mx-auto border border-pink-600/10">
          <h1 className="text-3xl font-extrabold mb-6 text-center text-pink-500">
            TurboGate
          </h1>

          <RouteForm onRouteAdded={fetchRoutes} />
          <RouteList routes={routes} />
        </div>
      </div>
    </>
  );
}

export default App;
