import { test, expect } from '@playwright/test';

test.describe('Content Viewing Strategies', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to courses and pick the first one
    await page.goto('http://localhost:5173/courses');
    await page.locator('.rt-Card').first().click();
  });

  test('should open video content in a focused screen', async ({ page }) => {
    // Look for a card that mentions 'video' (from badge or summary)
    const videoItem = page.locator('.rt-Card', { hasText: 'video' }).first();
    
    // We expect it to have a PlayCircle icon or just be clickable
    await videoItem.click();

    // Check URL has shifted to the focused viewer pattern
    await expect(page).toHaveURL(/\/courses\/.+\/content\/.+/);
    
    // Should see the theater mode button or similar focused UI
    await expect(page.getByRole('button', { name: /Theater Mode/i })).toBeVisible();
    await expect(page.locator('iframe')).toBeVisible();
  });

  test('should open PDF content in a dialog overlay', async ({ page }) => {
    // Look for a card that mentions 'pdf'
    const pdfItem = page.locator('.rt-Card', { hasText: 'pdf' }).first();
    
    await pdfItem.click();

    // URL should NOT change because it is a dialog
    await expect(page).not.toHaveURL(/\/courses\/.+\/content\/.+/);

    // Dialog should be visible
    const dialog = page.getByRole('dialog');
    await expect(dialog).toBeVisible();
    await expect(dialog.locator('iframe')).toBeVisible();
    
    // Can close it
    await dialog.getByRole('button', { name: /Close/i }).click();
    await expect(dialog).not.toBeVisible();
  });
});
