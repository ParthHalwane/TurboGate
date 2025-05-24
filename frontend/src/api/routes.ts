export type Route = {
    path: string;
    target: string;
  };
  
  const API_BASE = "http://localhost:10000"; // Replace with your backend URL
  
  export async function fetchRoutes(): Promise<Route[]> {
    const res = await fetch(`${API_BASE}/routes`);
    return res.json();
  }
  
  export async function addRoute(route: Route) {
    const res = await fetch(`${API_BASE}/routes`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(route),
    });
    return res.json();
  }
  
  export async function deleteRoute(path: string) {
    const res = await fetch(`${API_BASE}/routes/${encodeURIComponent(path)}`, {
      method: "DELETE",
    });
    return res.json();
  }
  