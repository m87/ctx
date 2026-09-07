import { TestBed } from '@angular/core/testing';
import { QueryFormComponent } from './query-form.component';

describe('QueryFormComponent', () => {
  it('offers saving when the current query can be saved', () => {
    TestBed.configureTestingModule({ imports: [QueryFormComponent] });
    const fixture = TestBed.createComponent(QueryFormComponent);
    fixture.componentRef.setInput('query', 'tag = focus');
    fixture.componentRef.setInput('canSave', true);
    let requestedSave = false;
    fixture.componentInstance.saveStart.subscribe(() => (requestedSave = true));
    fixture.detectChanges();

    const saveButton = Array.from(
      fixture.nativeElement.querySelectorAll('button') as NodeListOf<HTMLButtonElement>,
    ).find((button) => button.textContent?.includes('Save query'));
    saveButton?.click();

    expect(saveButton?.disabled).toBe(false);
    expect(requestedSave).toBe(true);
  });

  it('confirms a named query from the expanded save form', () => {
    TestBed.configureTestingModule({ imports: [QueryFormComponent] });
    const fixture = TestBed.createComponent(QueryFormComponent);
    fixture.componentRef.setInput('saveExpanded', true);
    fixture.componentRef.setInput('saveName', 'Focus');
    let confirmed = false;
    fixture.componentInstance.saveConfirm.subscribe(() => (confirmed = true));
    fixture.detectChanges();

    const nameInput = fixture.nativeElement.querySelector('#saved-query-name') as HTMLInputElement;
    nameInput.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true }));

    expect(confirmed).toBe(true);
  });
});
