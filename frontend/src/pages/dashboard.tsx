import { NextPage } from 'next';
import Head from 'next/head';
import { useAuth, withAuth } from '@/contexts/AuthContext';

const DashboardPage: NextPage = () => {
  const { user, logout, isAdmin } = useAuth();

  const handleLogout = () => {
    logout();
    window.location.href = '/login';
  };

  return (
    <>
      <Head>
        <title>Dashboard - CertiTrack</title>
        <meta name="description" content="CertiTrack Dashboard" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <div className="min-h-screen bg-gray-50">
        {/* Header */}
        <header className="bg-white shadow">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between items-center py-6">
              <div className="flex items-center">
                <h1 className="text-3xl font-bold text-gray-900">
                  Certi<span className="text-blue-600">Track</span>
                </h1>
              </div>
              <div className="flex items-center space-x-4">
                <div className="text-sm text-gray-700">
                  Welcome, <span className="font-medium">{user?.firstName} {user?.lastName}</span>
                  {isAdmin && (
                    <span className="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                      Admin
                    </span>
                  )}
                </div>
                <button
                  onClick={handleLogout}
                  className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors"
                >
                  Sign Out
                </button>
              </div>
            </div>
          </div>
        </header>

        {/* Main Content */}
        <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
          <div className="px-4 py-6 sm:px-0">
            {/* Welcome Section */}
            <div className="bg-white overflow-hidden shadow rounded-lg mb-6">
              <div className="px-4 py-5 sm:p-6">
                <h2 className="text-lg font-medium text-gray-900 mb-4">
                  Welcome to CertiTrack Dashboard
                </h2>
                <p className="text-gray-600 mb-4">
                  You have successfully logged in! This is a protected page that requires authentication.
                </p>

                {/* User Info */}
                <div className="bg-gray-50 rounded-lg p-4">
                  <h3 className="text-sm font-medium text-gray-900 mb-3">Your Account Information:</h3>
                  <dl className="grid grid-cols-1 gap-x-4 gap-y-3 sm:grid-cols-2">
                    <div>
                      <dt className="text-sm font-medium text-gray-500">Full Name</dt>
                      <dd className="text-sm text-gray-900">{user?.firstName} {user?.lastName}</dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">Email</dt>
                      <dd className="text-sm text-gray-900">{user?.email}</dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">Role</dt>
                      <dd className="text-sm text-gray-900 capitalize">{user?.role}</dd>
                    </div>
                    <div>
                      <dt className="text-sm font-medium text-gray-500">Status</dt>
                      <dd className="text-sm text-gray-900">
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${user?.isActive
                            ? 'bg-green-100 text-green-800'
                            : 'bg-red-100 text-red-800'
                          }`}>
                          {user?.isActive ? 'Active' : 'Inactive'}
                        </span>
                      </dd>
                    </div>
                    {user?.phone && (
                      <div>
                        <dt className="text-sm font-medium text-gray-500">Phone</dt>
                        <dd className="text-sm text-gray-900">{user.phone}</dd>
                      </div>
                    )}
                    <div>
                      <dt className="text-sm font-medium text-gray-500">Last Login</dt>
                      <dd className="text-sm text-gray-900">
                        {user?.lastLogin
                          ? new Date(user.lastLogin).toLocaleString()
                          : 'Never'
                        }
                      </dd>
                    </div>
                  </dl>
                </div>
              </div>
            </div>

            {/* Quick Stats */}
            <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-6">
              <div className="bg-white overflow-hidden shadow rounded-lg">
                <div className="p-5">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      <div className="w-8 h-8 bg-blue-500 rounded-md flex items-center justify-center">
                        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                        </svg>
                      </div>
                    </div>
                    <div className="ml-5 w-0 flex-1">
                      <dl>
                        <dt className="text-sm font-medium text-gray-500 truncate">Total People</dt>
                        <dd className="text-lg font-medium text-gray-900">0</dd>
                      </dl>
                    </div>
                  </div>
                </div>
              </div>

              <div className="bg-white overflow-hidden shadow rounded-lg">
                <div className="p-5">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      <div className="w-8 h-8 bg-green-500 rounded-md flex items-center justify-center">
                        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      </div>
                    </div>
                    <div className="ml-5 w-0 flex-1">
                      <dl>
                        <dt className="text-sm font-medium text-gray-500 truncate">Active Certifications</dt>
                        <dd className="text-lg font-medium text-gray-900">0</dd>
                      </dl>
                    </div>
                  </div>
                </div>
              </div>

              <div className="bg-white overflow-hidden shadow rounded-lg">
                <div className="p-5">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      <div className="w-8 h-8 bg-yellow-500 rounded-md flex items-center justify-center">
                        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
                        </svg>
                      </div>
                    </div>
                    <div className="ml-5 w-0 flex-1">
                      <dl>
                        <dt className="text-sm font-medium text-gray-500 truncate">Expiring Soon</dt>
                        <dd className="text-lg font-medium text-gray-900">0</dd>
                      </dl>
                    </div>
                  </div>
                </div>
              </div>

              <div className="bg-white overflow-hidden shadow rounded-lg">
                <div className="p-5">
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      <div className="w-8 h-8 bg-purple-500 rounded-md flex items-center justify-center">
                        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 7.172V5L8 4z" />
                        </svg>
                      </div>
                    </div>
                    <div className="ml-5 w-0 flex-1">
                      <dl>
                        <dt className="text-sm font-medium text-gray-500 truncate">Equipment</dt>
                        <dd className="text-lg font-medium text-gray-900">0</dd>
                      </dl>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Coming Soon Section */}
            <div className="bg-white overflow-hidden shadow rounded-lg">
              <div className="px-4 py-5 sm:p-6">
                <h3 className="text-lg font-medium text-gray-900 mb-4">Coming Soon</h3>
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
                  <div className="border border-gray-200 rounded-lg p-4">
                    <h4 className="font-medium text-gray-900">People Management</h4>
                    <p className="text-sm text-gray-600 mt-1">Add and manage personnel records</p>
                  </div>
                  <div className="border border-gray-200 rounded-lg p-4">
                    <h4 className="font-medium text-gray-900">Equipment Tracking</h4>
                    <p className="text-sm text-gray-600 mt-1">Track equipment and assets</p>
                  </div>
                  <div className="border border-gray-200 rounded-lg p-4">
                    <h4 className="font-medium text-gray-900">Certification Management</h4>
                    <p className="text-sm text-gray-600 mt-1">Manage certifications and renewals</p>
                  </div>
                  <div className="border border-gray-200 rounded-lg p-4">
                    <h4 className="font-medium text-gray-900">Automated Notifications</h4>
                    <p className="text-sm text-gray-600 mt-1">Email alerts for expiring certifications</p>
                  </div>
                  <div className="border border-gray-200 rounded-lg p-4">
                    <h4 className="font-medium text-gray-900">Reports & Analytics</h4>
                    <p className="text-sm text-gray-600 mt-1">Generate detailed compliance reports</p>
                  </div>
                  <div className="border border-gray-200 rounded-lg p-4">
                    <h4 className="font-medium text-gray-900">Document Management</h4>
                    <p className="text-sm text-gray-600 mt-1">Upload and manage certification documents</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </main>
      </div>
    </>
  );
};

export default withAuth(DashboardPage);