import { useState } from 'react';
import { MapComponent } from './components/MapComponent';
import { SuggestionList } from './components/SuggestionList';
import './index.css';

function App() {
  const [weather, setWeather] = useState(null);
  const [suggestions, setSuggestions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);

  const handleLocationSelect = async (lat: number, lng: number) => {
    setLoading(true);
    setIsSidebarOpen(true);
    console.log("Sending request to backend:", { lat, lng });
    try {
      const response = await fetch('http://localhost:8080/api/suggest', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ lat, lng }),
      });

      console.log("Backend response status:", response.status);

      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      const data = await response.json();
      console.log("Backend data received:", data);
      setWeather(data.weather);
      setSuggestions(data.suggestions);
    } catch (error) {
      console.error('Error fetching suggestions:', error);
      alert('Failed to fetch suggestions. Ensure backend is running.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative w-screen h-screen overflow-hidden font-sans text-gray-900">
      {/* Map Layer */}
      <div className="absolute inset-0 z-0">
        <MapComponent onLocationSelect={handleLocationSelect} />
      </div>

      {/* Sidebar Overlay */}
      <div
        className={`absolute top-0 left-0 h-full w-full md:w-96 z-10 transition-transform duration-300 transform ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full'
          } pointer-events-none`}
      >
        <div className="h-full pointer-events-auto">
          <SuggestionList weather={weather} suggestions={suggestions} loading={loading} />
        </div>

        {/* Toggle Button */}
        <button
          onClick={() => setIsSidebarOpen(!isSidebarOpen)}
          className="absolute top-1/2 -right-8 w-8 h-16 bg-white rounded-r-xl shadow-md flex items-center justify-center text-gray-600 hover:text-blue-500 focus:outline-none pointer-events-auto"
          title="Toggle Sidebar"
        >
          {isSidebarOpen ? '❮' : '❯'}
        </button>
      </div>

      <div className="absolute top-4 right-4 z-20 pointer-events-none">
        <h1 className="text-2xl font-bold text-white drop-shadow-md bg-black/20 backdrop-blur-sm px-4 py-2 rounded-lg">
          Weather & Map Planner
        </h1>
      </div>
    </div>
  );
}

export default App;
