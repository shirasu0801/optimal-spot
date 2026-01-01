import React from 'react';

interface WeatherInfo {
    main: string;
    description: string;
    temp: number;
}

interface Spot {
    id: string;
    name: string;
    latitude: number;
    longitude: number;
    rating: number;
    user_ratings_total: number;
    types: string[];
    photo_reference: string;
    crowd_level: string;
    weather_suitability: string;
    score: number;
}

interface SuggestionListProps {
    weather: WeatherInfo | null;
    suggestions: Spot[];
    loading: boolean;
}

export const SuggestionList: React.FC<SuggestionListProps> = ({ weather, suggestions, loading }) => {
    if (loading) {
        return <div className="p-4 text-center">Loading suggestions...</div>;
    }

    if (!weather && suggestions.length === 0) {
        return (
            <div className="p-8 text-center text-gray-500">
                <h2 className="text-xl font-bold mb-2">Welcome!</h2>
                <p>Drop a pin on the map to get sightseeing suggestions.</p>
            </div>
        );
    }

    return (
        <div className="h-full overflow-y-auto bg-white/90 backdrop-blur-md shadow-xl p-6 rounded-r-2xl border-r border-white/20">
            {weather && (
                <div className="mb-6 p-4 bg-blue-50 rounded-xl border border-blue-100">
                    <h3 className="text-sm font-semibold text-blue-800 uppercase tracking-wider mb-2">Current Weather</h3>
                    <div className="flex items-center justify-between">
                        <div>
                            <span className="text-3xl font-bold text-gray-800">{Math.round(weather.temp)}°C</span>
                            <p className="text-gray-600 capitalize">{weather.description}</p>
                        </div>
                        <div className="text-right">
                            <span className="px-3 py-1 bg-white rounded-full text-xs font-medium text-blue-600 shadow-sm">
                                {weather.main}
                            </span>
                        </div>
                    </div>
                </div>
            )}

            <h3 className="text-lg font-bold text-gray-800 mb-4 border-b pb-2">Top 3 Suggestions</h3>

            <div className="space-y-4">
                {suggestions.map((spot, index) => (
                    <div key={spot.id} className="bg-white rounded-xl p-4 shadow-sm hover:shadow-md transition-shadow border border-gray-100">
                        <div className="flex items-start justify-between mb-2">
                            <h4 className="font-bold text-gray-900 text-lg">
                                <span className="text-blue-500 mr-2">#{index + 1}</span>
                                {spot.name}
                            </h4>
                            <span className="bg-yellow-100 text-yellow-800 text-xs px-2 py-1 rounded-full flex items-center">
                                ★ {spot.rating}
                            </span>
                        </div>

                        <div className="flex flex-wrap gap-2 mb-3">
                            <span className={`text-xs px-2 py-1 rounded-full ${spot.weather_suitability.includes("Good") ? 'bg-green-100 text-green-700' : 'bg-orange-100 text-orange-700'
                                }`}>
                                Weather: {spot.weather_suitability}
                            </span>
                            <span className={`text-xs px-2 py-1 rounded-full ${spot.crowd_level === 'Low' ? 'bg-green-100 text-green-700' :
                                    spot.crowd_level === 'Medium' ? 'bg-yellow-100 text-yellow-700' : 'bg-red-100 text-red-700'
                                }`}>
                                Crowd: {spot.crowd_level}
                            </span>
                        </div>

                        <div className="text-xs text-gray-500">
                            {spot.types.slice(0, 3).join(", ")} • {spot.user_ratings_total} reviews
                        </div>
                    </div>
                ))}
            </div>

            {suggestions.length === 0 && weather && (
                <p className="text-center text-gray-500 mt-8">No spots found nearby.</p>
            )}
        </div>
    );
};
