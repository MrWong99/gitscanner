import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ScanOverviewComponent } from './scan-overview.component';

describe('ScanOverviewComponent', () => {
  let component: ScanOverviewComponent;
  let fixture: ComponentFixture<ScanOverviewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ScanOverviewComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ScanOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
