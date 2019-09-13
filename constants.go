//
//  OI Commands Quick Reference
//
// =====================================================================================================================
//  ---------------------  ----------  --------------------------------------------------------------------------
//   Command                Opcode      [ DataByte1, DataByte2, etc. ]
//  ---------------------  ----------  --------------------------------------------------------------------------
//   Start                  128         []
//   Baud                   129         [ BaudCode ]
//   Control                130         []
//   Safe                   131         []
//   Full                   132         []
//   Power                  133         []
//   Spot                   134         []
//   Clean                  135         []
//   Max Clean              136         []
//   Drive                  137         [ VelocityHigh, VelocityLow, RadiusHigh, RadiusLow ]
//   Motors                 138         [ MotorsState ]
//   Leds                   139         [ LedsState, PowerColor, PowerIntensity ]
//   Song                   140         [ SongNum, SongLength ]
//   Play                   141         [ SongNum ]
//   Query                  142         [ Packet ]
//   Force Seeking Dock     143         []
//   Pwm Motors             144         [ MainBrushPwm, SideBrushPwm, VacuumPwm ]
//   Drive Wheels           145         [ RightVelocityHigh, RightVelocityLow, LeftVelocityHigh, LeftVelocityLow ]
//   Drive Pwm              146         [ RightPwmHigh, RightPwmLow, LeftPwmHigh, LeftPwmLow ]
//   Stream                 148         [ NumPackets ]
//   Query List             149         [ NumPackets ]
//   Do Stream              150         [ StreamState ]
//   Scheduling Leds        162         [ weekdays, SchedulingLedsState ]
//   Digit Leds Raw         163         [ Digit3, Digit2, Digit1, Digit0 ]
//   Digit Leds Ascii       164         [ Digit3, Digit2, Digit1, Digit0 ]
//   Buttons                165         [ Buttons ]
//   Schedule               167         [ Days, SunHour, SunMin, MonHour, etc. ]
//   Set Day/Time           168         [ Day, Hour, Minute ]
//   Stop                   173         []
//  ---------------------  ----------  --------------------------------------------------------------------------
// =====================================================================================================================
//

package oi-bot
