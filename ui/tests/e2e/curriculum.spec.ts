import { test, expect } from '@playwright/test';

test.describe('Teacher Curriculum Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:5173/courses');
    
    // Switch to Instructor role to see management buttons
    await page.getByRole('button', { name: 'Learner' }).click();
    await page.getByRole('menuitem').filter({ hasText: 'Instructor' }).click();
  });

  test('should create a new course', async ({ page }) => {
    const uniqueTitle = `Automated Test Course ${Date.now()}`;
    
    // Click Create New Course
    await page.getByRole('button', { name: /Create New Course/i }).click();

    // Fill out the form
    await page.getByPlaceholder(/Advanced Go Patterns/i).fill(uniqueTitle);
    await page.getByPlaceholder(/e.g. Engineering/i).fill('Testing');
    await page.getByPlaceholder(/Briefly describe/i).fill('This is an automated test course summary.');
    
    // Submit
    await page.getByRole('button', { name: /Create Course/i }).click();

    // Should navigate to detail page
    await expect(page).toHaveURL(/\/courses\/.+/);
    await expect(page.getByRole('heading', { name: uniqueTitle })).toBeVisible();
  });

  test('should add a section to a course', async ({ page }) => {
    // Navigate to the first existing course
    await page.locator('a:has-text("Courses")').click();
    await page.locator('.rt-Card').first().click();

    const sectionTitle = `Test Section ${Date.now()}`;
    
    // Click Add Section
    await page.getByRole('button', { name: /Add Section/i }).click();
    
    // Fill form
    await page.getByPlaceholder(/e.g. Getting Started/i).fill(sectionTitle);
    
    // Submit
    await page.getByRole('button', { name: 'Add Section' }).click();

    // Should appear in the list
    await expect(page.getByText(sectionTitle)).toBeVisible();
  });
});
