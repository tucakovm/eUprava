import { inject } from '@angular/core';
import { CanActivateFn, CanMatchFn, Router, UrlSegment, ActivatedRouteSnapshot, Route } from '@angular/router';
import { AuthService } from './services/auth.service';

function checkAccess(allowed: string[] | undefined): boolean {
  const auth = inject(AuthService);
  const router = inject(Router);

  const role = auth.userRole;
  if (!auth.isAuthenticated || !role) {
    router.navigate(['/login']);
    return false;
  }

  // ako ruta nema definisane role, tretiraj kao zabranjeno
  if (!allowed || allowed.length === 0) {
    router.navigate(['/unathorized']); // ili /login
    return false;
  }

  if (allowed.includes(role)) {
    return true;
  }

  router.navigate(['/unathorized']); 
  return false;
}

export const roleCanActivate: CanActivateFn = (route: ActivatedRouteSnapshot) => {
  const roles = route.data['roles'] as string[] | undefined;
  return checkAccess(roles);
};

export const roleCanMatch: CanMatchFn = (route: Route, segments: UrlSegment[]) => {
  const roles = (route.data?.['roles'] as string[] | undefined) ?? undefined;
  return checkAccess(roles);
};
