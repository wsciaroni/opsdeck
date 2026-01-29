import { useAuth } from '../context/AuthContext';
import { BarChart } from 'lucide-react';
import EmptyState from '../components/EmptyState';

export default function Reports() {
    const { currentOrg } = useAuth();

    if (!currentOrg) return <div>Please select an organization</div>;

    return (
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-semibold text-gray-900">Reports</h1>
            </div>

            <EmptyState
                title="No reports available"
                description="Reporting features are coming soon."
                icon={BarChart}
            />
        </div>
    );
}
