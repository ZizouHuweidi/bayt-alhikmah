import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Library, Loader2, Mail, Shield, User } from "lucide-react";
import { useEffect } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { useAuth } from "@/lib/auth/AuthContext";

export const Route = createFileRoute("/settings")({
	component: SettingsPage,
});

function SettingsPage() {
	const navigate = useNavigate();
	const { isAuthenticated, isLoading, user } = useAuth();

	useEffect(() => {
		if (!isLoading && !isAuthenticated) {
			navigate({ to: "/login" });
		}
	}, [isAuthenticated, isLoading, navigate]);

	if (isLoading) {
		return (
			<div className="flex min-h-screen items-center justify-center bg-slate-50">
				<Loader2 className="h-8 w-8 animate-spin text-emerald-600" />
			</div>
		);
	}

	return (
		<div className="min-h-screen bg-slate-50 py-8">
			<div className="mx-auto max-w-2xl px-4">
				<div className="mb-8 flex items-center gap-3">
					<div className="rounded-xl bg-emerald-600 p-2">
						<Library className="h-6 w-6 text-white" />
					</div>
					<h1 className="text-2xl font-bold text-slate-900">
						Account Settings
					</h1>
				</div>

				<Card>
					<CardContent className="space-y-6 pt-6">
						<div className="flex items-center gap-4">
							<div className="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100">
								<User className="h-6 w-6 text-emerald-600" />
							</div>
							<div>
								<p className="font-medium text-slate-900">{user.username}</p>
								<p className="text-sm text-slate-500">Username</p>
							</div>
						</div>

						<div className="flex items-center gap-4">
							<div className="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100">
								<Mail className="h-6 w-6 text-emerald-600" />
							</div>
							<div>
								<p className="font-medium text-slate-900">
									{user.email || "No email"}
								</p>
								<p className="text-sm text-slate-500">Email address</p>
							</div>
						</div>

						<div className="flex items-center gap-4">
							<div className="flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100">
								<Shield className="h-6 w-6 text-emerald-600" />
							</div>
							<div>
								<p className="font-medium text-slate-900">
									Password auth enabled
								</p>
								<p className="text-sm text-slate-500">
									OAuth support can be added later.
								</p>
							</div>
						</div>
					</CardContent>
				</Card>
			</div>
		</div>
	);
}
