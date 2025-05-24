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
      `http://localhost:10000/api/add-route?${query.toString()}`
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
        <label className="block font-medium">Route Path</label>
        <input
          type="text"
          value={route}
          onChange={(e) => setRoute(e.target.value)}
          className="w-full px-4 py-2 border rounded-md"
          placeholder="/new"
          required
        />
      </div>
      <div>
        <label className="block font-medium">Domain URL</label>
        <input
          type="url"
          value={domain}
          onChange={(e) => setDomain(e.target.value)}
          className="w-full px-4 py-2 border rounded-md"
          placeholder="https://example.com"
          required
        />
      </div>
      <button
        type="submit"
        className="w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700"
      >
        ➕ Add Route
      </button>
    </form>
  );
}

export default RouteForm;
