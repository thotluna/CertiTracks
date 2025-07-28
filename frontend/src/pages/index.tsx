import { NextPage } from 'next';
import Head from 'next/head';
import Link from 'next/link';

const HomePage: NextPage = () => {
  return (
    <>
      <Head>
        <title>CertiTrack - Certification Management System</title>
        <meta name="description" content="Manage certifications for people and equipment" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
        <div className="container mx-auto px-4 py-16">
          {/* Header */}
          <div className="text-center mb-16">
            <h1 className="text-5xl font-bold text-gray-900 mb-4">
              Certi<span className="text-primary-600">Track</span>
            </h1>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Comprehensive certification management system for people and equipment.
              Never miss an expiration date again.
            </p>
          </div>

          {/* Features Grid */}
          <div className="grid md:grid-cols-3 gap-8 mb-16">
            <div className="bg-white rounded-lg shadow-soft p-6 text-center">
              <div className="w-12 h-12 bg-primary-100 rounded-lg flex items-center justify-center mx-auto mb-4">
                <svg className="w-6 h-6 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Track Certifications</h3>
              <p className="text-gray-600">
                Manage all certifications for personnel and equipment in one centralized system.
              </p>
            </div>

            <div className="bg-white rounded-lg shadow-soft p-6 text-center">
              <div className="w-12 h-12 bg-warning-100 rounded-lg flex items-center justify-center mx-auto mb-4">
                <svg className="w-6 h-6 text-warning-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Smart Alerts</h3>
              <p className="text-gray-600">
                Receive automated email notifications before certifications expire.
              </p>
            </div>

            <div className="bg-white rounded-lg shadow-soft p-6 text-center">
              <div className="w-12 h-12 bg-success-100 rounded-lg flex items-center justify-center mx-auto mb-4">
                <svg className="w-6 h-6 text-success-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Detailed Reports</h3>
              <p className="text-gray-600">
                Generate comprehensive reports and analytics for compliance tracking.
              </p>
            </div>
          </div>

          {/* CTA Section */}
          <div className="text-center">
            <div className="bg-white rounded-lg shadow-medium p-8 max-w-md mx-auto">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Get Started</h2>
              <p className="text-gray-600 mb-6">
                Ready to streamline your certification management?
              </p>
              <div className="space-y-3">
                <Link
                  href="/login"
                  className="block w-full bg-primary-600 text-white py-3 px-6 rounded-lg font-medium hover:bg-primary-700 transition-colors"
                >
                  Sign In
                </Link>
                <Link
                  href="/register"
                  className="block w-full border border-gray-300 text-gray-700 py-3 px-6 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                >
                  Create Account
                </Link>
              </div>
            </div>
          </div>

          {/* Status Section */}
          <div className="mt-16 text-center">
            <div className="inline-flex items-center space-x-2 bg-white rounded-full px-4 py-2 shadow-soft">
              <div className="w-2 h-2 bg-success-500 rounded-full animate-pulse"></div>
              <span className="text-sm text-gray-600">System Status: Online</span>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default HomePage;