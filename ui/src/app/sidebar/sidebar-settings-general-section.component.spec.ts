import { TestBed } from '@angular/core/testing';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { SettingsMutations } from '../../api/settings.mutations';
import { settingsQueryKeys, SettingsQueries } from '../../api/settings.queries';
import { timeZoneSettingKey, TimeZoneService } from '../shared/time-zone.service';
import { SidebarSettingsGeneralSectionComponent } from './sidebar-settings-general-section.component';

describe('SidebarSettingsGeneralSectionComponent', () => {
  it('selects the saved time zone when settings are already loaded', async () => {
    const settings = {
      'client.general.theme': 'light',
      'client.general.firstDay': 'Monday',
      [timeZoneSettingKey]: 'Europe/Warsaw',
    };
    const queryClient = new QueryClient({
      defaultOptions: { queries: { staleTime: Infinity } },
    });
    queryClient.setQueryData(settingsQueryKeys.settings(), settings);

    await TestBed.configureTestingModule({
      imports: [SidebarSettingsGeneralSectionComponent],
      providers: [
        provideTanStackQuery(queryClient),
        {
          provide: SettingsQueries,
          useValue: {
            settings: () => ({
              queryKey: settingsQueryKeys.settings(),
              queryFn: async () => settings,
              staleTime: Infinity,
            }),
          },
        },
        {
          provide: SettingsMutations,
          useValue: {
            save: () => ({ mutationFn: async () => undefined }),
          },
        },
        {
          provide: TimeZoneService,
          useValue: {
            options: ['UTC', 'Europe/Warsaw'],
            setPreference: vi.fn(),
          },
        },
      ],
    }).compileComponents();

    const fixture = TestBed.createComponent(SidebarSettingsGeneralSectionComponent);
    fixture.detectChanges();
    await fixture.whenStable();
    fixture.detectChanges();

    const select = fixture.nativeElement.querySelector('#general-time-zone') as HTMLSelectElement;
    expect(select?.value).toBe('Europe/Warsaw');
  });
});
