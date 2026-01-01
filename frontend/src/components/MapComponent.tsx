import React, { useState } from 'react';
import { APIProvider, Map, Marker } from '@vis.gl/react-google-maps';

interface MapComponentProps {
  onLocationSelect: (lat: number, lng: number) => void;
}

export const MapComponent: React.FC<MapComponentProps> = ({ onLocationSelect }) => {
  const [markerPosition, setMarkerPosition] = useState<{ lat: number; lng: number } | null>(null);

  const API_KEY = import.meta.env.VITE_GOOGLE_MAPS_API_KEY || '';

  const handleMapClick = (ev: any) => {
    console.log("Raw Map Click Event:", ev);

    let lat, lng;

    // Try to extract latLng from different possible event structures
    if (ev.detail && ev.detail.latLng) {
      lat = ev.detail.latLng.lat;
      lng = ev.detail.latLng.lng;
    } else if (ev.latLng) {
      lat = ev.latLng.lat(); // Google Maps API latLng object uses functions
      lng = ev.latLng.lng();
    }

    if (lat !== undefined && lng !== undefined) {
      console.log("Coordinates extracted successfully:", lat, lng);
      setMarkerPosition({ lat, lng });
      onLocationSelect(lat, lng);
    } else {
      console.error("Could not extract coordinates from click event. Please check the 'Raw Map Click Event' log above.");
      // Fallback for debugging: click center if extraction fails? No, that's confusing.
    }
  };

  return (
    <div style={{ height: '100vh', width: '100%' }}>
      <APIProvider apiKey={API_KEY}>
        <Map
          style={{ width: '100%', height: '100%' }}
          defaultCenter={{ lat: 35.6895, lng: 139.6917 }} // Tokyo
          defaultZoom={13}
          gestureHandling={'greedy'}
          disableDefaultUI={false}
          onClick={handleMapClick}
          mapId="DEMO_MAP_ID" // Required for Advanced Markers if needed
        >
          {markerPosition && (
            <Marker position={markerPosition} />
          )}
        </Map>
      </APIProvider>
    </div>
  );
};
