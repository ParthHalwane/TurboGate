export type Route = {
  path: string;
  target: string;
};

interface Props {
  routes: Route[];
}

function RouteList({ routes }: Props) {
  return (
    <div>
      <h2 className="text-xl text-center font-semibold mb-2">📃 Current Routes</h2>
      <p className="text-sm text-gray-400 mb-4">
        ⚠️ Note: Some websites may not load correctly due to HTTPS restrictions
        or browser security policies (like CORS or certificate validation). Unfortunately, we can not bypass them.
      </p>

      <ul className="list-disc list-inside space-y-2">
        {routes.map((r, idx) => (
          <li key={idx}>
            <span className="font-mono text-blue-400">{r.path}</span> →
            <a
              href={r.target}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-200 underline hover:text-blue-100 ml-1"
            >
              {r.target}
            </a>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default RouteList;
