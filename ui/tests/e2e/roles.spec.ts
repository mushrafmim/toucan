import { test, expect } from '@playwright/test';

test.describe('Role Switching and Navigation', () => {
  test.beforeEach(async ({ page }) => {
    // Standard Vite port
    await page.goto('http://localhost:5173/');
  });

  test('should default to Learner view and show basic nav', async ({ page }) => {
    // Check role switcher shows Student (id: learner)
    await expect(page.getByRole('button', { name: 'Student' })).toBeVisible();

    const sidebar = page.getByRole('navigation', { name: /Primary/i });
    await expect(sidebar.getByRole('link', { name: /Overview/i })).toBeVisible();
    await expect(sidebar.getByRole('link', { name: /Courses/i })).toBeVisible();
    
    // Admin/Teacher specific items should be hidden
    await expect(sidebar.getByRole('link', { name: /Curriculum/i })).not.toBeVisible();
    await expect(sidebar.getByRole('link', { name: /Users/i })).not.toBeVisible();
  });

  test('should switch to Instructor view and show curriculum', async ({ page }) => {
    // Open role switcher
    await page.getByRole('button', { name: 'Student' }).click();
    
    // Select Instructor
    await page.getByRole('menuitem').filter({ hasText: 'Instructor' }).click();

    // Trigger should update
    await expect(page.getByRole('button', { name: 'Instructor' })).toBeVisible();

    // Sidebar should update
    const sidebar = page.getByRole('navigation', { name: /Primary/i });
    await expect(sidebar.getByRole('link', { name: /Curriculum/i })).toBeVisible();
    await expect(sidebar.getByRole('link', { name: /Users/i })).not.toBeVisible();
  });

  test('should switch to Administrator view and show all items', async ({ page }) => {
    // Open role switcher
    await page.getByRole('button', { name: 'Student' }).click();
    
    // Select Administrator
    await page.getByRole('menuitem').filter({ hasText: 'Administrator' }).click();

    // Trigger should update
    await expect(page.getByRole('button', { name: 'Administrator' })).toBeVisible();

    // Sidebar should show Admin items
    const sidebar = page.getByRole('navigation', { name: /Primary/i });
    await expect(sidebar.getByRole('link', { name: /Users/i })).toBeVisible();
    await expect(sidebar.getByRole('link', { name: /Settings/i })).toBeVisible();
  });
});
