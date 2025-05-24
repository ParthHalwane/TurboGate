import { useState } from "react";
import toast from "react-hot-toast";

interface Props {
  onRouteAdded: () => void;
}

function RouteForm({ onRouteAdded }: Props) {
  const [route, setRoute] = useState("");
  const [domain, setDomain] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const query = new URLSearchParams({
      route: route.trim(),
      domain: domain.trim(),
    });

    const res = await fetch(
      `${import.meta.env.VITE_TURBOGATE_BACKEND_URL}/api/add-route?${query.toString()}`
    );

    if (res.ok) {
      toast.success(" Route added successfully!");
      setRoute("");
      setDomain("");
      onRouteAdded();
    } else {
      const err = await res.text();
      toast.error(`❌ Failed: ${err}`);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4 mb-6">
      <div>
        <label className="block font-medium text-pink-400">Route Path</label>
        <input
          type="text"
          value={route}
          onChange={(e) => setRoute(e.target.value)}
          className="w-full px-4 py-2 border border-pink-200 rounded-md"
          placeholder="/new"
          required
        />
      </div>
      <div>
        <label className="block font-medium text-pink-400">Domain URL</label>
        <input
          type="url"
          value={domain}
          onChange={(e) => setDomain(e.target.value)}
          className="w-full px-4 py-2 border border-pink-200 rounded-md"
          placeholder="https://example.com"
          required
        />
      </div>
      <button
        type="submit"
        className="w-full bg-pink-500 text-white py-2 rounded-md hover:bg-pink-600"
      >
        ➕ Add Route
      </button>
    </form>
  );
}

export default RouteForm;
