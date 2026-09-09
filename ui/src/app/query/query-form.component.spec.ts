import { TestBed } from '@angular/core/testing';
import {
  getQuerySuggestionContext,
  QUERY_SUGGESTIONS,
  QueryFormComponent,
} from './query-form.component';

describe('query suggestions', () => {
  it('contains all supported query terms', () => {
    expect(QUERY_SUGGESTIONS).toEqual([
      'and',
      'or',
      'in',
      '!=',
      '==',
      '~',
      '<',
      '>',
      '<=',
      '>=',
      'text',
      'name',
      'label',
      'project',
      'start',
      'end',
      'duration',
      'sessions',
    ]);
  });

  it('filters and locates the term under the cursor', () => {
    expect(getQuerySuggestionContext('project == na', 13)).toEqual({
      start: 11,
      end: 13,
      suggestions: ['name'],
    });
    expect(getQuerySuggestionContext('sessions !', 10).suggestions).toEqual(['!=']);
  });

  it('does not offer query terms inside quoted text', () => {
    expect(getQuerySuggestionContext('name == "pro', 12).suggestions).toEqual([]);
  });
});

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

  it('filters suggestions and inserts the active one with Enter', () => {
    TestBed.configureTestingModule({ imports: [QueryFormComponent] });
    const fixture = TestBed.createComponent(QueryFormComponent);
    const queryChanges: string[] = [];
    fixture.componentInstance.queryChange.subscribe((query) => queryChanges.push(query));
    fixture.detectChanges();

    const textarea = fixture.nativeElement.querySelector('#context-query') as HTMLTextAreaElement;
    textarea.value = 'na';
    textarea.setSelectionRange(2, 2);
    textarea.dispatchEvent(new Event('input', { bubbles: true }));
    fixture.detectChanges();

    const options = fixture.nativeElement.querySelectorAll(
      '#context-query-suggestions [role="option"]',
    ) as NodeListOf<HTMLButtonElement>;
    expect(Array.from(options).map((option) => option.textContent?.trim())).toEqual(['name']);

    textarea.dispatchEvent(new KeyboardEvent('keydown', { key: 'Enter', bubbles: true }));
    fixture.detectChanges();

    expect(textarea.value).toBe('name');
    expect(queryChanges).toEqual(['na', 'name']);
    expect(fixture.nativeElement.querySelector('#context-query-suggestions')).toBeNull();
  });

  it('inserts symbolic operators with Tab', () => {
    TestBed.configureTestingModule({ imports: [QueryFormComponent] });
    const fixture = TestBed.createComponent(QueryFormComponent);
    fixture.detectChanges();

    const textarea = fixture.nativeElement.querySelector('#context-query') as HTMLTextAreaElement;
    textarea.value = 'name !';
    textarea.setSelectionRange(6, 6);
    textarea.dispatchEvent(new Event('input', { bubbles: true }));
    fixture.detectChanges();
    textarea.dispatchEvent(new KeyboardEvent('keydown', { key: 'Tab', bubbles: true }));

    expect(textarea.value).toBe('name !=');
  });
});
