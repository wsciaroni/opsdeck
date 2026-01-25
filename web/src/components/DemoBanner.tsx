import { Info } from 'lucide-react';

export default function DemoBanner() {
  return (
    <div className="bg-blue-600 text-white px-4 py-2 text-center text-sm font-medium">
      <div className="flex items-center justify-center gap-2">
        <Info className="h-4 w-4" />
        <p>This is a demo instance. Data may be deleted at any time.</p>
      </div>
    </div>
  );
}
