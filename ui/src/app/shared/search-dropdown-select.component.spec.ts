import { TestBed } from '@angular/core/testing';
import { SearchDropdownSelectComponent } from './search-dropdown-select.component';

describe('SearchDropdownSelectComponent', () => {
  it('uses the compact trigger by default', () => {
    const fixture = createFixture();
    const trigger = fixture.nativeElement.querySelector(
      'button[aria-haspopup="listbox"]',
    ) as HTMLButtonElement;

    expect(trigger.classList.contains('h-8')).toBe(true);
    expect(trigger.textContent).not.toContain('Workspace');
  });

  it('renders the verbose variant with secondary information', () => {
    const fixture = createFixture('verbose');
    const trigger = fixture.nativeElement.querySelector(
      'button[aria-haspopup="listbox"]',
    ) as HTMLButtonElement;

    expect(trigger.classList.contains('min-h-12')).toBe(true);
    expect(trigger.textContent).toContain('Personal');
    expect(trigger.textContent).toContain('Private time tracking');
  });
});

function createFixture(variant: 'default' | 'verbose' = 'default') {
  TestBed.configureTestingModule({ imports: [SearchDropdownSelectComponent] });
  const fixture = TestBed.createComponent(SearchDropdownSelectComponent);
  fixture.componentRef.setInput('inputId', 'workspace-select');
  fixture.componentRef.setInput('variant', variant);
  fixture.componentRef.setInput('options', [
    {
      value: 'workspace-1',
      label: 'Personal',
      description: 'Private time tracking',
      color: '#6366f1',
    },
  ]);
  fixture.componentRef.setInput('value', 'workspace-1');
  fixture.detectChanges();
  return fixture;
}
